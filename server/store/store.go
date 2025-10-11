package store

import (
	"fmt"
	"getMeMod/server/store/core"
	"getMeMod/server/store/utils"
	"getMeMod/utils/logger"
	"sync"
	"time"
)

type Store struct {
	mu                      sync.RWMutex
	basePath                string
	hashTable               *core.HashTable
	segmentManager          *core.SegmentManager
	compactionResultChannel chan *core.CompactionResult
	doneChannel             chan struct{}
}

func NewStore(mainBasePath, compactedBasePath string) *Store {

	// unbuffered channel to ensure that compaction results are processed in order
	compactionResultChannel := make(chan *core.CompactionResult, 0)

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
	s.mu.RLock()
	defer s.mu.RUnlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hashTable.Size()
}

func (s *Store) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashTable.Clear()
	if err := s.segmentManager.Clear(); err != nil {
		return fmt.Errorf("failed to clear segment manager: %w", err)
	}
	return nil
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hashTable.Keys()
}

func (s *Store) BatchSet(batch map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(batch) == 0 {
		return nil
	}

	logger.Info("Starting BatchSet operation for", len(batch), "items")

	// In-memory buffer to hold serialized entries
	writeBuffer := make([]byte, 0, 64*1024) // Start with 64KB capacity
	// Map to hold entries for the current chunk
	chunkEntries := make(map[string]*core.Entry)

	for key, value := range batch {
		keyBytes := s.convertStringToBytes(key)
		valueBytes := s.convertStringToBytes(value)
		timeStamp := time.Now().UnixNano()

		entry, err := core.CreateEntry(keyBytes, valueBytes, timeStamp)
		if err != nil {
			// This should ideally not happen.
			logger.Error("BatchSet: Failed to create entry for key", key, ":", err)
			continue
		}

		serializedEntry, err := entry.Serialize()
		if err != nil {
			logger.Error("BatchSet: Failed to serialize entry for key", key, ":", err)
			continue
		}

		// If adding the new entry exceeds the buffer, flush the current buffer first.
		if len(writeBuffer)+len(serializedEntry) > 64*1024 && len(writeBuffer) > 0 {
			logger.Info("BatchSet: Write buffer full, flushing", len(chunkEntries), "entries.")
			newIndexPointers, err := s.segmentManager.FlushBuffer(writeBuffer, chunkEntries)
			if err != nil {
				// This is a critical error. The data is partially written to disk but the index is not updated.
				// This could lead to data inconsistency. Future improvements could involve a recovery mechanism.
				return fmt.Errorf("failed to flush write buffer during batch set: %w", err)
			}
			s.hashTable.BatchUpdate(newIndexPointers)

			// Reset buffers for the next chunk
			writeBuffer = make([]byte, 0, 64*1024)
			chunkEntries = make(map[string]*core.Entry)
		}

		writeBuffer = append(writeBuffer, serializedEntry...)
		chunkEntries[key] = entry
	}

	// Flush any remaining entries in the buffer
	if len(writeBuffer) > 0 {
		logger.Info("BatchSet: Flushing remaining", len(chunkEntries), "entries.")
		newIndexPointers, err := s.segmentManager.FlushBuffer(writeBuffer, chunkEntries)
		if err != nil {
			return fmt.Errorf("failed to flush final write buffer during batch set: %w", err)
		}
		s.hashTable.BatchUpdate(newIndexPointers)
	}

	logger.Info("BatchSet operation completed successfully.")
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.compactionResultChannel)

	close(s.doneChannel)

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
