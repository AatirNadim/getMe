package store

import (
	"fmt"
	"getMeMod/server/store/core"
	"getMeMod/server/store/utils"
	"getMeMod/server/store/utils/constants"
	"getMeMod/server/utils/logger"
	"time"
)

type Store struct {
	// mu                      sync.RWMutex
	basePath                string
	hashTable               *core.HashTable
	segmentManager          *core.SegmentManager
	compactionResultChannel chan *core.CompactionResult
	doneChannel             chan struct{}
}

func NewStore(mainBasePath, compactedBasePath string) *Store {

	// unbuffered channel to ensure that compaction results are processed in order
	compactionResultChannel := make(chan *core.CompactionResult)

	hashTable := core.NewHashTable()
	segmentManager, err := core.NewSegmentManager(mainBasePath, compactedBasePath, hashTable, compactionResultChannel)
	if err != nil {
		panic(err)
	}

	logger.Info("creating a new store instance on the base path:", mainBasePath)

	store := &Store{
		basePath:                mainBasePath,
		hashTable:               hashTable,
		segmentManager:          segmentManager,
		compactionResultChannel: compactionResultChannel,
		doneChannel:             make(chan struct{}),
	}

	go store.listenForCompactionResults()

	return store
}

func (s *Store) Get(key string) (string, bool, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	logger.Info("Getting the file and the offset for key:", key)
	hashTableEntry, exists := s.hashTable.Get(key)
	if !exists {
		logger.Error("key not found: ", key)
		return "", false, utils.ErrKeyNotFound
	}

	data, _, err := s.segmentManager.Read(hashTableEntry.SegmentId, hashTableEntry.Offset)
	if err != nil {
		return "", false, err
	}

	return s.convertBytesToString(data.Value), true, nil
}

func (s *Store) Put(key string, value string) error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	logger.Info("Putting key:", key, "with value: ", value)

	keyBytes := s.convertStringToBytes(key)
	valueBytes := s.convertStringToBytes(value)

	timeStamp := time.Now().UnixNano()

	entry, err := core.CreateEntry(keyBytes, valueBytes, timeStamp)
	if err != nil {
		return err
	}

	logger.Info("appending entry with key:", key, " to segment manager")

	segmentId, offset, err := s.segmentManager.Append(entry)
	if err != nil {
		return err
	}

	logger.Info("updating hash table with key:", key, " segmentId:", segmentId, " offset:", offset)

	s.hashTable.Put(key, segmentId, offset, timeStamp, entry.ValueSize)

	// if newSegmentCreated {
	// 	go s.performCompaction()
	// }

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
	writeBuffer := make([]byte, 0, constants.MaxChunkSize) // Start with MaxChunkSize capacity
	// Map to hold entries for the current chunk
	chunkEntries := make([]*core.Entry, 0, len(batch))

	flushAndReset := func() error {
		logger.Debug("flushAndReset called with buffer size:", len(writeBuffer), "and", len(chunkEntries), "entries")
		if len(writeBuffer) == 0 {
			return nil
		}
		logger.Info("BatchPut: Flushing buffer with", len(chunkEntries), "entries.")
		flushResults, err := s.segmentManager.FlushBuffer(writeBuffer, chunkEntries)
		if err != nil {
			return fmt.Errorf("failed to flush write buffer during batch set: %w", err)
		}

		newIndexPointers := make(map[string]*core.HashTableEntry)
		for i, result := range flushResults {
			originalEntry := chunkEntries[i]
			// Map the original entry key to its new location
			newIndexPointers[s.convertBytesToString(originalEntry.Key)] = &core.HashTableEntry{
				SegmentId: uint32(result.SegmentID),
				Offset:    uint32(result.Offset),
				TimeStamp: originalEntry.TimeStamp,
				ValueSize: originalEntry.ValueSize,
			}
		}

		s.hashTable.BatchUpdate(newIndexPointers)

		// Reset buffers for the next chunk
		writeBuffer = make([]byte, 0, 64*1024)
		chunkEntries = make([]*core.Entry, 0, len(batch))
		return nil
	}

	for key, value := range batch {
		keyBytes := s.convertStringToBytes(key)
		valueBytes := s.convertStringToBytes(value)
		timeStamp := time.Now().UnixNano()

		entry, err := core.CreateEntry(keyBytes, valueBytes, timeStamp)
		if err != nil {
			// This should ideally not happen.
			logger.Error("BatchPut: Failed to create entry for key", key, ":", err)
			continue
		}

		serializedEntry, err := entry.Serialize()
		if err != nil {
			logger.Error("BatchPut: Failed to serialize entry for key", key, ":", err)
			continue
		}

		// If adding the new entry exceeds the buffer, flush the current buffer first.
		if len(writeBuffer)+len(serializedEntry) > constants.MaxChunkSize {
			logger.Debug("BatchPut: Buffer full. Flushing current buffer before adding key:", key)
			if err := flushAndReset(); err != nil {
				return err
			}
		}

		writeBuffer = append(writeBuffer, serializedEntry...)
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

func (s *Store) Close() error {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	close(s.doneChannel)

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
	for {
		select {
		case compactionResult := <-s.compactionResultChannel:
			logger.Info("Received compaction result with", compactionResult.CompactedHashTable.Size(), "entries and", len(compactionResult.OldSegmentIds), "old segments to delete")
			s.applyCompactionResult(compactionResult)
		case <-s.doneChannel:
			logger.Info("Shutting down compaction result listener")
			return
		}
	}
}

// func (s *Store) performCompaction() {

// 	if totalSegments > constants.ThresholdForCompaction {
// 		logger.Info("total segments:", totalSegments, "exceeds threshold:", constants.ThresholdForCompaction, "starting compaction")

// 		s.segmentManager.PerformCompaction(s.hashTable, s.compactedSegmentManager)
// 	}
// }
