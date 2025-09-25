package store

import (
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
	segmentManager, err := NewSegmentManager(basePath)

	if err != nil {
		panic(err)
	}

	return &Store{
		basePath:      basePath,
		hashTable:     NewHashTable(),
		segmentManager: segmentManager,
	}
}

func (s *Store) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.hashTable.Get(key)
	if !exists {
		return "", false, nil
	}

	data, err := s.segmentManager.Read(entry.segmentId, entry.offset)
	if err != nil {
		return "", false, err
	}

	return s.convertBytesToString(data.Value), true, nil
}

func (s *Store) Put(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyBytes := s.convertStringToBytes(key)
	valueBytes := s.convertStringToBytes(value)

	entry := &Entry{
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(keyBytes)),
		ValueSize: uint32(len(valueBytes)),
		Key:       keyBytes,
		Value:     valueBytes,
	}

	segmentId, offset, err := s.segmentManager.Append(entry)
	if err != nil {
		return err
	}

	s.hashTable.Put(key, segmentId, offset)
	return nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.segmentManager.Delete(s.convertStringToBytes(key))
	if err != nil {
		return err
	}
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