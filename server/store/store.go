package store

import (
	"bytes"
	"fmt"
	"maps"
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

func (s *Store) BatchGet(keys []string) (serverUtils.BatchGetResult, error) {
	// 1. Lock & Lookup
	hashTableEntries, notFoundList := s.hashTable.GetBatch(keys)

	// 2. Group and Sort by Location
	type BatchEntry struct {
		Key   string
		Entry *core.HashTableEntry
	}

	segmentMap := make(map[uint32][]*BatchEntry) // map from segmentId to list of entries in that segment
	for key, entry := range hashTableEntries {
		segmentMap[entry.SegmentId] = append(segmentMap[entry.SegmentId], &BatchEntry{Key: key, Entry: entry})
	}

	// sort the entries for each segment by offset to optimize disk reads (sequential reads are faster)
	for _, entries := range segmentMap {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Entry.Offset < entries[j].Entry.Offset
		})
	}

	// 3. Perform Efficient Disk Reads in Parallel
	var wg sync.WaitGroup
	resultsChan := make(chan serverUtils.BatchGetResult, len(segmentMap))

	for segmentId, entries := range segmentMap {
		wg.Add(1)
		go func(segId uint32, batchEntries []*BatchEntry) {
			defer wg.Done()

			localResult := serverUtils.BatchGetResult{
				Found:    make(map[string]string),
				NotFound: make([]string, 0),
				Errors:   make(map[string]string),
			}

			offsets := make([]uint32, len(batchEntries))
			segmentKeys := make([]string, len(batchEntries))
			for i, be := range batchEntries {
				offsets[i] = be.Entry.Offset
				segmentKeys[i] = be.Key
			}

			// load the actual segment keys as well, so that in case of any error, we can update the final result with the corresponding keys and their error messages
			readEntries, errorsMap, err := s.segmentManager.BatchRead(segId, offsets, segmentKeys)

			if err != nil {
				if err == utils.ErrSegmentNotFound {
					// Fallback: Try to get keys individually (might be compacted)
					for _, be := range batchEntries {
						val, found, err2 := s.Get(be.Key)
						if err2 == nil && found {
							localResult.Found[be.Key] = val
						} else {
							if err2 != nil {
								localResult.Errors[be.Key] = err2.Error()
							} else {
								localResult.NotFound = append(localResult.NotFound, be.Key)
							}
						}
					}
					resultsChan <- localResult
					return
				}
				logger.Error("Error batch reading segment:", segId, "error:", err)
				// in case of an error different than segment not found, we consider all the keys in this batch as failed with the same error message, since we cannot determine which key caused the error
				for _, be := range batchEntries {
					localResult.Errors[be.Key] = err.Error()
				}
				resultsChan <- localResult
				return
			}

			for k, v := range errorsMap {
				localResult.Errors[k] = v
			}

			// Add found entries
			for _, e := range readEntries {
				if e != nil {
					if e.ValueSize > 0 {
						localResult.Found[string(e.Key)] = s.convertBytesToString(e.Value)
					} else {
						localResult.NotFound = append(localResult.NotFound, string(e.Key))
					}
				}
			}

			logger.Info("local result for the segment: \n", segmentId, "\nis\n", localResult)

			resultsChan <- localResult
		}(segmentId, entries)
	}

	wg.Wait()
	close(resultsChan)

	// 5. Assemble and Return
	finalResult := serverUtils.BatchGetResult{
		Found:    make(map[string]string),
		NotFound: notFoundList,
		Errors:   make(map[string]string),
	}

	for res := range resultsChan {
		maps.Copy(finalResult.Found, res.Found)
		finalResult.NotFound = append(finalResult.NotFound, res.NotFound...)
		maps.Copy(finalResult.Errors, res.Errors)
	}

	return finalResult, nil
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

func (s *Store) BatchPut(batch map[string]string) (serverUtils.BatchPutResult, error) {
	result := serverUtils.BatchPutResult{
		Successful: 0,
		Failed:     make(map[string]string),
	}

	if len(batch) == 0 {
		return result, nil
	}
	// No top-level lock here to allow for concurrent reads.
	// Locking is handled at a granular level in SegmentManager and HashTable.

	logger.Info("Starting BatchPut operation for", len(batch), "items")

	// In-memory buffer to hold serialized entries

	// Note: we are using a buffer pool to fetch a buffer from, which reduces the number of allocations and GC overhead for large batches, since we reuse the same buffer for multiple batch put operations
	writeBuffer := batchBufferPool.Get().(*bytes.Buffer)
	defer batchBufferPool.Put(writeBuffer)
	writeBuffer.Reset()

	// Start with MaxChunkSize capacity
	// Map to hold entries for the current chunk
	chunkEntries := make([]*core.Entry, 0, len(batch))
	newIndexPointers := make(map[string]*core.HashTableEntry)

	flushAndReset := func() {
		// logger.Debug("flushAndReset called with buffer size:", writeBuffer.Len(), "and", len(chunkEntries), "entries")

		if writeBuffer.Len() == 0 {
			return
		}
		// logger.Info("BatchPut: Flushing buffer with", len(chunkEntries), "entries.")
		flushResults, err := s.segmentManager.FlushBuffer(writeBuffer.Bytes(), chunkEntries)
		if err != nil {
			errStr := fmt.Sprintf("failed to flush write buffer: %v", err)
			logger.Error("BatchPut: " + errStr)
			for _, entry := range chunkEntries {
				keyStr := s.convertBytesToString(entry.Key)
				result.Failed[keyStr] = errStr
				entryPool.Put(entry)
			}
		} else {
			for i, res := range flushResults {
				originalEntry := chunkEntries[i]
				// Map the original entry key to its new location

				indexEntry := &core.HashTableEntry{
					SegmentId: uint32(res.SegmentID),
					Offset:    uint32(res.Offset),
					TimeStamp: originalEntry.TimeStamp,
					ValueSize: originalEntry.ValueSize,
				}
				keyStr := s.convertBytesToString(originalEntry.Key)
				newIndexPointers[keyStr] = indexEntry
			}

			// logger.Debug("batchput: updating hashtable with the latest index pointers, --> ", newIndexPointers)

			s.hashTable.BatchUpdate(newIndexPointers)

			// release all the entry objects back into the pool
			for _, entry := range chunkEntries {
				entryPool.Put(entry)
			}

			result.Successful += len(chunkEntries)
		}

		// Reset buffers for the next chunk
		writeBuffer.Reset()
		chunkEntries = chunkEntries[:0]
		clear(newIndexPointers)
	}

	// Process each key-value pair in the batch
	for key, value := range batch {
		keyBuffer := keyBufferPool.Get().(*bytes.Buffer) // fetch a buffer from the pool
		keyBuffer.Reset()
		keyBuffer.WriteString(key)

		valueBuffer := valueBufferPool.Get().(*bytes.Buffer) // fetch a buffer from the pool
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
			errStr := fmt.Sprintf("failed to serialize entry: %v", err)
			logger.Error("BatchPut: Failed to serialize entry for key", key, ":", err)
			result.Failed[key] = errStr
			// in case of any problem, release the object back into the pool
			entryPool.Put(entry)
			continue
		}

		// If adding the new entry exceeds the buffer, flush the current buffer first.
		if writeBuffer.Len()+len(serializedEntry) > constants.MaxChunkSize {
			// logger.Debug("BatchPut: Buffer full. Flushing current buffer before adding key:", key)
			flushAndReset()
		}

		writeBuffer.Write(serializedEntry)
		// chunkEntries have the order of entries same as they are added to the writeBuffer
		chunkEntries = append(chunkEntries, entry)
	}

	// Flush any remaining entries in the buffer
	flushAndReset()

	logger.Info("BatchPut operation completed. Successful:", result.Successful, "Failed:", len(result.Failed))
	return result, nil
}

func (s *Store) BatchDelete(keys []string) (serverUtils.BatchDeleteResult, error) {
	result := serverUtils.BatchDeleteResult{
		Successful: 0,
		Failed:     make(map[string]string),
	}

	if len(keys) == 0 {
		return result, nil
	}
	logger.Info("Starting BatchDelete operation for", len(keys), "keys")

	writeBuffer := batchBufferPool.Get().(*bytes.Buffer)
	defer batchBufferPool.Put(writeBuffer)
	writeBuffer.Reset()

	chunkEntries := make([]*core.Entry, 0, len(keys))

	flushAndReset := func() {
		if writeBuffer.Len() == 0 {
			return
		}
		_, err := s.segmentManager.FlushBuffer(writeBuffer.Bytes(), chunkEntries)
		if err != nil {
			errStr := fmt.Sprintf("failed to flush write buffer during batch delete: %v", err)
			logger.Error(errStr)
			for _, entry := range chunkEntries {
				keyStr := s.convertBytesToString(entry.Key)
				result.Failed[keyStr] = errStr
				entryPool.Put(entry)
			}
		} else {
			keysToDelete := make([]string, 0, len(chunkEntries))
			for _, entry := range chunkEntries {
				keysToDelete = append(keysToDelete, string(entry.Key))
				entryPool.Put(entry)
			}

			s.hashTable.BatchDelete(keysToDelete)
			result.Successful += len(chunkEntries)
		}

		writeBuffer.Reset()
		chunkEntries = chunkEntries[:0]
	}

	for _, key := range keys {
		// Check existence to avoid writing useless tombstones
		if _, exists := s.hashTable.Get(key); !exists {
			// Idempotent delete: if it doesn't exist, it's effectively "deleted"
			result.Successful++
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
			errStr := fmt.Sprintf("failed to create deletion entry: %v", err)
			logger.Error("BatchDelete: Failed to create deletion entry for key", key, ":", err)
			result.Failed[key] = errStr
			continue
		}

		serializedEntry, err := entry.Serialize()
		if err != nil {
			errStr := fmt.Sprintf("failed to serialize entry: %v", err)
			logger.Error("BatchDelete: Failed to serialize entry for key", key, ":", err)
			result.Failed[key] = errStr
			entryPool.Put(entry)
			continue
		}

		if writeBuffer.Len()+len(serializedEntry) > constants.MaxChunkSize {
			flushAndReset()
		}

		writeBuffer.Write(serializedEntry)
		chunkEntries = append(chunkEntries, entry)
	}

	flushAndReset()

	logger.Info("BatchDelete operation completed successfully.")
	return result, nil
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
