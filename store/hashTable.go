package store

import "sync"


type HashTableEntry struct {
	segmentId	uint32
	offset		uint32
}


type HashTable struct {
	mu sync.RWMutex
	table map[string]HashTableEntry
}

func NewHashTable() *HashTable {
	return &HashTable{
		table: make(map[string]HashTableEntry),
	}
}

func (ht *HashTable) Get(key string) (HashTableEntry, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	entry, exists := ht.table[key]
	return entry, exists
}

func (ht *HashTable) Put(key string, entry HashTableEntry) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.table[key] = entry
}

func (ht *HashTable) Delete(key string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	delete(ht.table, key)
}

func (ht *HashTable) Size() int {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	return len(ht.table)
}

func (ht *HashTable) Clear() {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.table = make(map[string]HashTableEntry)
}

func (ht *HashTable) Keys() []string {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	keys := make([]string, 0, len(ht.table))
	for k := range ht.table {
		keys = append(keys, k)
	}
	return keys
}

func (ht *HashTable) Entries() map[string]HashTableEntry {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	entries := make(map[string]HashTableEntry, len(ht.table))
	for k, v := range ht.table {
		entries[k] = v
	}
	return entries
}