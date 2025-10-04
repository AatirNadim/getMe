package core

import "sync"

type HashTableEntry struct {
	TimeStamp int64
	SegmentId uint32
	Offset    uint32
	ValueSize uint32
}

// a hash table does not to be concerned about how the raw data is stored in the disk, it only deals with the mapping from key to (segmentId, offset)
type HashTable struct {
	mu    sync.RWMutex
	table map[string]*HashTableEntry
}

func NewHashTable() *HashTable {
	return &HashTable{
		table: make(map[string]*HashTableEntry),
	}
}

// func (ht *HashTable) IsEntryPresentInHashTable(key string) bool {
// 	ht.mu.RLock()
// 	defer ht.mu.RUnlock()
// 	_, exists := ht.Get(key)
// 	return exists
// }

func (ht *HashTable) Get(key string) (*HashTableEntry, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	entry, exists := ht.table[key]
	return entry, exists
}

// if the existing entry's timestamp is greater than or equal to the new entry's timestamp, it means the new entry is older or same, so we do not consider it present
func (ht *HashTable) Put(key string, segmentId uint32, offset uint32, timeStamp int64, valueSize uint32) error {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	// this is not going to fail even if the key does not exist in the hashtable
	if existingEntry, ok := ht.table[key]; ok {
		if timeStamp == existingEntry.TimeStamp && offset < existingEntry.Offset {
			return nil // incoming entry timestamp is same but offset is less, which means its an older entry, do nothing
		}
		if timeStamp < existingEntry.TimeStamp {
			return nil // incoming entry is older, do nothing
		}
	}

	ht.table[key] = &HashTableEntry{
		SegmentId: segmentId,
		Offset:    offset,
		TimeStamp: timeStamp,
		ValueSize: valueSize,
	}
	return nil
}

func (ht *HashTable) PutEntry(key string, entry HashTableEntry) {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	ht.table[key] = &entry
}

func (ht *HashTable) Merge(other *HashTable) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	for key, otherEntry := range other.table {
		if existingEntry, ok := ht.table[key]; ok {
			// If the other entry is newer, update the table
			if otherEntry.TimeStamp > existingEntry.TimeStamp {
				ht.table[key] = otherEntry
			}
			// If timestamps are equal, keep the one with the higher offset
			if otherEntry.TimeStamp == existingEntry.TimeStamp && otherEntry.Offset > existingEntry.Offset {
				ht.table[key] = otherEntry
			}
		} else {
			// If the key doesn't exist, just add it
			ht.table[key] = otherEntry
		}
	}
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
	ht.table = make(map[string]*HashTableEntry)
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
		entries[k] = *v
	}
	return entries
}

func (ht *HashTable) DeleteDeletionEntries() {
	ht.mu.Lock()
	defer ht.mu.Unlock()
	for key, entry := range ht.table {
		if entry.ValueSize == 0 {
			delete(ht.table, key)
		}
	}
}

