package store

import (
	"encoding/binary"
	"time"
	"getMeMod/store/utils"
	"getMeMod/store/logger"
)


// we are dealing with the segment ids instead of the actual segment locations


// the size of an entry instance will be variable, depending on the key and value sizes
type Entry struct {
	TimeStamp uint32
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}


func CreateDeletionEntry(key []byte) (*Entry, error) {
	return &Entry {
		TimeStamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: 0,
		Key:       key,
		Value:     nil,
	}, nil
}


func (e *Entry) IsDeletionEntry() bool {
	return e.ValueSize == 0
}

func (e *Entry) getEntryKVPairSize() uint32 {
	return e.KeySize + e.ValueSize;
}


func (e *Entry) getEntrySize() uint32 {
	return e.getEntryKVPairSize() + 12; // 12 bytes for the headers
}



func (e *Entry) Serialize() ([]byte, error) {
	logger.Info("Serializing entry with key: ", string(e.Key), " and value size: ", e.ValueSize)

	bytarr := make([]byte, e.getEntrySize())

	offset := 0


	binary.LittleEndian.PutUint32(bytarr[offset:], e.TimeStamp)

	offset += 4

	binary.LittleEndian.PutUint32(bytarr[offset:], e.KeySize)

	offset += 4


	binary.LittleEndian.PutUint32(bytarr[offset:], e.ValueSize)

	offset += 4

	copy(bytarr[offset:], e.Key)

	offset += int(e.KeySize)

	if(e.ValueSize > 0) {
		copy(bytarr[offset:], e.Value)
	}

	return bytarr, nil

}


func DeserializeEntry(bytarr []byte) (*Entry, error) {

	logger.Info("Deserializing entry")

	offset := 0

	e := &Entry{}
	e.TimeStamp = binary.LittleEndian.Uint32(bytarr[offset:])

	offset += 4

	e.KeySize = binary.LittleEndian.Uint32(bytarr[offset:])

	offset += 4

	e.ValueSize = binary.LittleEndian.Uint32(bytarr[offset:])

	if int(e.KeySize)+int(e.ValueSize) != len(bytarr)-12 {
		return nil, utils.ErrInvalidEntry
	}

	offset += 4

	e.Key = make([]byte, e.KeySize)
	copy(e.Key, bytarr[offset:offset+int(e.KeySize)])

	offset += int(e.KeySize)

	if(e.ValueSize > 0) {
		e.Value = make([]byte, e.ValueSize)
		copy(e.Value, bytarr[offset:offset+int(e.ValueSize)])
	} else {
		e.Value = nil
	}

	return e, nil
}



func (entry *Entry) getHashTableFromEntry() *HashTableEntry {
	return &HashTableEntry{
		segmentId: 0, // will be set later
		offset:    0, // will be set later
		timeStamp: entry.TimeStamp,
	}
}