package store

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/AatirNadim/getMe/server/store/core"
	"github.com/AatirNadim/getMe/server/store/utils"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
	serverUtils "github.com/AatirNadim/getMe/server/utils"
	"github.com/AatirNadim/getMe/server/utils/logger"
)

// Pools for reusing byte buffers to reduce allocations for keys and values.
var (
	keyBufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	valueBufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	batchBufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	entryPool = sync.Pool{
		New: func() interface{} {
			return new(core.Entry)
		},
	}
	hashTableEntryPool = sync.Pool{
		New: func() interface{} {
			return new(core.HashTableEntry)
		},
	}
)

type Store struct {
	// mu                      sync.RWMutex
	basePath                string
	hashTable               *core.HashTable
	segmentManager          *core.SegmentManager
	compactionResultChannel chan *core.CompactionResult
	doneChannel             chan struct{}
	wg                      sync.WaitGroup
}

func NewStore(mainBasePath, compactedBasePath string) *Store {

	// unbuffered channel to ensure that compaction results are processed in order
	compactionResultChannel := make(chan *core.CompactionResult)

	hashTable := core.NewHashTable()
	segmentManager, err := core.NewSegmentManager(mainBasePath, compactedBasePath, hashTable, compactionResultChannel)
	if err != nil {
		panic(err)
	}

	store := &Store{
		basePath:                mainBasePath,
		hashTable:               hashTable,
		segmentManager:          segmentManager,
		compactionResultChannel: compactionResultChannel,
		doneChannel:             make(chan struct{}),
	}

	store.wg.Add(1)
	go store.listenForCompactionResults()

	return store
}

func (s *Store) Get(key string) (string, bool, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	// logger.Info("Getting the file and the offset for key:", key)
	hashTableEntry, exists := s.hashTable.Get(key)
	if !exists {
		fmt.Println("key not found for the given request")
		return "", false, utils.ErrKeyNotFound
	}

	// fmt.Println("Getting the file and the offset for key:", key)
	// fmt.Println("store: 63: Found in hashtable with segmentId:", hashTableEntry.SegmentId, " offset:", hashTableEntry.Offset)

	data, _, err := s.segmentManager.Read(hashTableEntry.SegmentId, hashTableEntry.Offset)
	if err != nil {

		if err == utils.ErrSegmentNotFound {
			// logger.Error("segment not found for key:", key, " segmentId:", hashTableEntry.SegmentId)
			hashTableEntry, exists := s.hashTable.Get(key)
			if !exists {
				// logger.Error("key not found: ", key)
				return "", false, utils.ErrKeyNotFound
			}
			data, _, err = s.segmentManager.Read(hashTableEntry.SegmentId, hashTableEntry.Offset)
			if err != nil {
				// If it fails a second time, it's a real error.
				return "", false, err
			}
		} else {
			return "", false, err
		}

	}

	return s.convertBytesToString(data.Value), true, nil
}

func (s *Store) BatchGet(keys []string) ([]serverUtils.KeyValue, error) {
	// 1. Lock & Lookup
	hashTableEntries := s.hashTable.GetBatch(keys)

	// 2. Group and Sort by Location
	type BatchEntry struct {
		Key   string
		Entry *core.HashTableEntry
	}

	segmentMap := make(map[uint32][]*BatchEntry)
	for key, entry := range hashTableEntries {
		segmentMap[entry.SegmentId] = append(segmentMap[entry.SegmentId], &BatchEntry{Key: key, Entry: entry})
	}

	// Sort by Offset
	for _, entries := range segmentMap {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Entry.Offset < entries[j].Entry.Offset
		})
	}

	// 3. Release Lock (Already done in GetBatch)

	// 4. Perform Efficient Disk Reads in Parallel
	var wg sync.WaitGroup
	resultsChan := make(chan []*core.Entry, len(segmentMap))

	for segmentId, entries := range segmentMap {
		wg.Add(1)
		go func(segId uint32, batchEntries []*BatchEntry) {
			defer wg.Done()
			offsets := make([]uint32, len(batchEntries))
			for i, be := range batchEntries {
				offsets[i] = be.Entry.Offset
			}

			entries, err := s.segmentManager.BatchRead(segId, offsets)
			if err != nil {
				if err == utils.ErrSegmentNotFound {
					// Fallback: Try to get keys individually (might be compacted)
					fallbackEntries := make([]*core.Entry, 0)
					for _, be := range batchEntries {
						val, found, err2 := s.Get(be.Key)
						if err2 == nil && found {
							// Reconstruct entry (expensive but fallback)
							entry := &core.Entry{
								Key:       []byte(be.Key),
								Value:     []byte(val),
								ValueSize: uint32(len(val)),
							}
							fallbackEntries = append(fallbackEntries, entry)
						} else {
							logger.Error("Segment not found and key not found in fallback:", segId, "key:", be.Key)
						}
					}
					// Send empty or partial
					resultsChan <- fallbackEntries
					return
				}
				logger.Error("Error batch reading segment:", segId, "error:", err)
				resultsChan <- nil
				return
			}
			resultsChan <- entries
		}(segmentId, entries)
	}

	wg.Wait()
	close(resultsChan)

	// 5. Assemble and Return
	tempMap := make(map[string]string)
	for entries := range resultsChan {
		for _, e := range entries {
			if e != nil && e.ValueSize > 0 {
				tempMap[string(e.Key)] = s.convertBytesToString(e.Value)
			}
		}
	}

	// Sort in the order of requested keys
	orderedResult := make([]serverUtils.KeyValue, 0, len(keys))
	for _, k := range keys {
		if v, ok := tempMap[k]; ok {
			orderedResult = append(orderedResult, serverUtils.KeyValue{Key: k, Value: v})
		}
	}

	return orderedResult, nil
}

func (s *Store) Put(key string, value string) error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	logger.Info("Putting key:", key, "with value: ", value)

	// create a buffer for key and value to avoid multiple allocations
	keyBuffer := keyBufferPool.Get().(*bytes.Buffer)
	defer keyBufferPool.Put(keyBuffer)
	keyBuffer.Reset()
	keyBuffer.WriteString(key)

	valueBuffer := valueBufferPool.Get().(*bytes.Buffer)
	defer valueBufferPool.Put(valueBuffer)
	valueBuffer.Reset()
	valueBuffer.WriteString(value)

	// commented out to use the buffer pools instead

	// keyBytes := s.convertStringToBytes(key)
	// valueBytes := s.convertStringToBytes(value)

	timeStamp := time.Now().UnixNano()

	entry, err := core.CreateEntry(keyBuffer.Bytes(), valueBuffer.Bytes(), timeStamp)
	if err != nil {
		return err
	}

	// logger.Info("appending entry with key:", key, " to segment manager")

	segmentId, offset, err := s.segmentManager.Append(entry)
	if err != nil {
		return err
	}

	// logger.Info("updating hash table with key:", key, " segmentId:", segmentId, " offset:", offset)

	s.hashTable.Put(key, segmentId, offset, timeStamp, entry.ValueSize)
	logger.Info("key has been added and hashtable has been updated, key = ", key)

	return nil
}

func (s *Store) Delete(key string) error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	if _, exists := s.hashTable.Get(key); !exists {
		logger.Error("key not found: ", key)
		return utils.ErrKeyNotFound
	}

	logger.Info("creating deletion entry for key:", key)

	timeStamp := time.Now().UnixNano()

	deletionEntry, deletionEntryCreationErr := core.CreateDeletionEntry(s.convertStringToBytes(key), timeStamp)

	if deletionEntryCreationErr != nil {
		logger.Error("store: failed to create deletion entry:", deletionEntryCreationErr)
		return deletionEntryCreationErr
	}

	logger.Info("appending deletion entry for key:", key, " to segment manager")
	_, err := s.segmentManager.Delete(deletionEntry)
	if err != nil {
		return err
	}

	logger.Info("removing key:", key, " from hash table")
	s.hashTable.Delete(key)
	return nil
}

func (s *Store) Size() int {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	return s.hashTable.Size()
}

func (s *Store) Clear() error {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	s.hashTable.Clear()
	if err := s.segmentManager.Clear(); err != nil {
		return fmt.Errorf("failed to clear segment manager: %w", err)
	}
	return nil
}

func (s *Store) Keys() []string {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	return s.hashTable.Keys()
}

func (s *Store) BatchPut(batch map[string]string) error {
	if len(batch) == 0 {
		return nil
	}
	// No top-level lock here to allow for concurrent reads.
	// Locking is handled at a granular level in SegmentManager and HashTable.

	logger.Info("Starting BatchPut operation for", len(batch), "items")

	// In-memory buffer to hold serialized entries
	// writeBuffer := make([]byte, 0, constants.MaxChunkSize)
	writeBuffer := batchBufferPool.Get().(*bytes.Buffer)
	defer batchBufferPool.Put(writeBuffer)
	writeBuffer.Reset()

	// Start with MaxChunkSize capacity
	// Map to hold entries for the current chunk
	chunkEntries := make([]*core.Entry, 0, len(batch))
	newIndexPointers := make(map[string]*core.HashTableEntry)

	flushAndReset := func() error {
		// logger.Debug("flushAndReset called with buffer size:", writeBuffer.Len(), "and", len(chunkEntries), "entries")

		if writeBuffer.Len() == 0 {
			return nil
		}
		// logger.Info("BatchPut: Flushing buffer with", len(chunkEntries), "entries.")
		flushResults, err := s.segmentManager.FlushBuffer(writeBuffer.Bytes(), chunkEntries)
		if err != nil {
			return fmt.Errorf("failed to flush write buffer during batch set: %w", err)
		}

		for i, result := range flushResults {
			originalEntry := chunkEntries[i]
			// Map the original entry key to its new location

			indexEntry := hashTableEntryPool.Get().(*core.HashTableEntry)
			indexEntry.SegmentId = uint32(result.SegmentID)
			indexEntry.Offset = uint32(result.Offset)
			indexEntry.TimeStamp = originalEntry.TimeStamp
			indexEntry.ValueSize = originalEntry.ValueSize
			keyStr := s.convertBytesToString(originalEntry.Key)
			newIndexPointers[keyStr] = indexEntry
		}

		// logger.Debug("batchput: updating hashtable with the latest index pointers, --> ", newIndexPointers)

		s.hashTable.BatchUpdate(newIndexPointers)

		// release all the entry objects back into the pool
		for _, entry := range chunkEntries {
			entryPool.Put(entry)
		}

		// release all the hashTableEntry objects back into the pool
		for _, indexEntry := range newIndexPointers {
			hashTableEntryPool.Put(indexEntry)
		}

		// Reset buffers for the next chunk
		writeBuffer.Reset()
		chunkEntries = chunkEntries[:0]
		clear(newIndexPointers)
		return nil
	}

	// Process each key-value pair in the batch
	for key, value := range batch {
		keyBuffer := keyBufferPool.Get().(*bytes.Buffer)
		keyBuffer.Reset()
		keyBuffer.WriteString(key)

		valueBuffer := valueBufferPool.Get().(*bytes.Buffer)
		valueBuffer.Reset()
		valueBuffer.WriteString(value)

		keyBytes := make([]byte, keyBuffer.Len())
		copy(keyBytes, keyBuffer.Bytes())

		valueBytes := make([]byte, valueBuffer.Len())
		copy(valueBytes, valueBuffer.Bytes())

		timeStamp := time.Now().UnixNano()

		// commented out to use the buffer pools instead
		// entry, err := core.CreateEntry(keyBuffer.Bytes(), valueBuffer.Bytes(), timeStamp)

		entry := entryPool.Get().(*core.Entry)
		entry.Key = keyBytes // dont use the original variable's bytes here, since they are passed by reference
		entry.Value = valueBytes
		entry.KeySize = uint32(len(keyBytes))
		entry.ValueSize = uint32(len(valueBytes))
		entry.TimeStamp = timeStamp

		keyBufferPool.Put(keyBuffer)
		valueBufferPool.Put(valueBuffer)

		serializedEntry, err := entry.Serialize()
		if err != nil {
			logger.Error("BatchPut: Failed to serialize entry for key", key, ":", err)
			// in case of any problem, release the object back into the pool
			entryPool.Put(entry)
			continue
		}

		// If adding the new entry exceeds the buffer, flush the current buffer first.
		if writeBuffer.Len()+len(serializedEntry) > constants.MaxChunkSize {
			// logger.Debug("BatchPut: Buffer full. Flushing current buffer before adding key:", key)
			if err := flushAndReset(); err != nil {
				return err
			}
		}

		writeBuffer.Write(serializedEntry)
		// chunkEntries have the order of entries same as they are added to the writeBuffer
		chunkEntries = append(chunkEntries, entry)
	}

	// Flush any remaining entries in the buffer
	if err := flushAndReset(); err != nil {
		return err
	}

	logger.Info("BatchPut operation completed successfully.")
	return nil
}

func (s *Store) BatchDelete(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	logger.Info("Starting BatchDelete operation for", len(keys), "keys")

	writeBuffer := batchBufferPool.Get().(*bytes.Buffer)
	defer batchBufferPool.Put(writeBuffer)
	writeBuffer.Reset()

	chunkEntries := make([]*core.Entry, 0, len(keys))

	flushAndReset := func() error {
		if writeBuffer.Len() == 0 {
			return nil
		}
		_, err := s.segmentManager.FlushBuffer(writeBuffer.Bytes(), chunkEntries)
		if err != nil {
			return fmt.Errorf("failed to flush write buffer during batch delete: %w", err)
		}

		keysToDelete := make([]string, 0, len(chunkEntries))
		for _, entry := range chunkEntries {
			keysToDelete = append(keysToDelete, string(entry.Key))
			entryPool.Put(entry)
		}

		s.hashTable.BatchDelete(keysToDelete)

		writeBuffer.Reset()
		chunkEntries = chunkEntries[:0]
		return nil
	}

	for _, key := range keys {
		// Check existence to avoid writing useless tombstones
		if _, exists := s.hashTable.Get(key); !exists {
			continue
		}

		keyBuffer := keyBufferPool.Get().(*bytes.Buffer)
		keyBuffer.Reset()
		keyBuffer.WriteString(key)

		keyBytes := make([]byte, keyBuffer.Len())
		copy(keyBytes, keyBuffer.Bytes())
		keyBufferPool.Put(keyBuffer)

		timeStamp := time.Now().UnixNano()

		entry, err := core.CreateDeletionEntry(keyBytes, timeStamp)
		if err != nil {
			logger.Error("BatchDelete: Failed to create deletion entry for key", key, ":", err)
			continue
		}

		serializedEntry, err := entry.Serialize()
		if err != nil {
			logger.Error("BatchDelete: Failed to serialize entry for key", key, ":", err)
			entryPool.Put(entry)
			continue
		}

		if writeBuffer.Len()+len(serializedEntry) > constants.MaxChunkSize {
			if err := flushAndReset(); err != nil {
				return err
			}
		}

		writeBuffer.Write(serializedEntry)
		chunkEntries = append(chunkEntries, entry)
	}

	if err := flushAndReset(); err != nil {
		return err
	}

	logger.Info("BatchDelete operation completed successfully.")
	return nil
}

func (s *Store) Close() error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	close(s.doneChannel)

	s.wg.Wait()

	close(s.compactionResultChannel)
	return nil
}

// private methods

func (s *Store) convertStringToBytes(str string) []byte {
	return []byte(str)
}

func (s *Store) convertBytesToString(b []byte) string {
	return string(b)

}

func (s *Store) applyCompactionResult(compactionResult *core.CompactionResult) {

	logger.Debug("Applying compaction result with", compactionResult.CompactedHashTable, "entries and", len(compactionResult.OldSegmentIds), "old segments to delete")

	s.hashTable.Merge(compactionResult.CompactedHashTable)
	s.hashTable.DeleteDeletionEntries() // remove deletion entries from the hash table
	s.segmentManager.DeleteOldSegments(compactionResult.OldSegmentIds)
}

func (s *Store) listenForCompactionResults() {
	defer s.wg.Done()
	for {
		select {
		case compactionResult := <-s.compactionResultChannel:
			logger.Info("Received compaction signal, applying compaction result")
			// logger.Info("Received compaction result with", compactionResult.CompactedHashTable.Size(), "entries and", len(compactionResult.OldSegmentIds), "old segments to delete")
			s.applyCompactionResult(compactionResult)
		case <-s.doneChannel:
			logger.Info("Shutting down compaction result listener")
			return
		}
	}
}
