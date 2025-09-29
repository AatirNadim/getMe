package store

import (
	"getMeMod/store/core"
	"getMeMod/store/logger"
	"getMeMod/store/utils"
	"sync"
	"time"
)

type Store struct {
	mu             sync.RWMutex
	basePath       string
	hashTable      *core.HashTable
	segmentManager *core.SegmentManager
}

func NewStore(basePath string) *Store {
	hashTable := core.NewHashTable()
	segmentManager, err := core.NewSegmentManager(basePath, hashTable)

	if err != nil {
		panic(err)
	}

	logger.Info("creating a new store instance on the base path:", basePath)

	return &Store{
		basePath:       basePath,
		hashTable:      hashTable,
		segmentManager: segmentManager,
	}
}

func (s *Store) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logger.Info("Getting the file and the offset for key:", key)
	entry, exists := s.hashTable.Get(key)
	if !exists {
		logger.Error("key not found: ", key)
		return "", false, utils.ErrKeyNotFound
	}

	data, _, err := s.segmentManager.Read(entry.SegmentId, entry.Offset)
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

	nextSegmentId, offset, err := s.segmentManager.Append(entry)
	if err != nil {
		return err
	}

	logger.Info("updating hash table with key:", key, " segmentId:", nextSegmentId-1, " offset:", offset)

	s.hashTable.Put(key, nextSegmentId-1, offset, timeStamp, entry.ValueSize)
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

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashTable.Clear()
	s.segmentManager.Clear()
}

func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hashTable.Keys()
}

// private methods

func (s *Store) convertStringToBytes(str string) []byte {
	return []byte(str)
}

func (s *Store) convertBytesToString(b []byte) string {
	return string(b)
}
