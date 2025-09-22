package store

import (

	"time"
)


// we are dealing with the segment ids instead of the actual segment locations


type Entry struct {
	TimeStamp uint32
	KeySize   uint16
	ValueSize uint16
	Key       []byte
	Value     []byte
}


func (e *Entry) addPreviousEntryDeletionEntry() *Entry {
	return &Entry {
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   e.KeySize,
		ValueSize: 0,
		Key:       e.Key,
		Value:     nil,
	}
}


func isDeletionEntry(e *Entry) bool {
	return e.ValueSize == 0
}