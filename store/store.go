package store

import (
	"getMeMod/store/logger"
	"getMeMod/store/utils"
	"sync"
	"time"
)

type Store struct {
	mu sync.RWMutex
	basePath    string
	hashTable   *HashTable
	segmentManager *SegmentManager
}


func NewStore(basePath string) *Store {
	hashTable := NewHashTable()
	segmentManager, err := NewSegmentManager(basePath, hashTable)

	if err != nil {
		panic(err)
	}

	logger.Info("creating a new store instance on the base path:", basePath)

	return &Store{
		basePath:      basePath,
		hashTable:     hashTable,
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

	data, _, err := s.segmentManager.Read(entry.segmentId, entry.offset)
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

	timeStamp := uint32(time.Now().Unix())

	entry, err := CreateEntry(keyBytes, valueBytes, timeStamp)
	if err != nil {
		return err
	}

	logger.Info("appending entry with key:", key, " to segment manager")

	nextSegmentId, offset, err := s.segmentManager.Append(entry)
	if err != nil {
		return err
	}

	logger.Info("updating hash table with key:", key, " segmentId:", nextSegmentId - 1, " offset:", offset)

	s.hashTable.Put(key, nextSegmentId - 1, offset, timeStamp, entry.ValueSize)
	return nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Deleting key:", key)

	_, err := s.segmentManager.Delete(s.convertStringToBytes(key))
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