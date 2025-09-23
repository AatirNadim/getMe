package store

import (
	"encoding/binary"
	"time"
	"getMeMod/store/utils"
)


// we are dealing with the segment ids instead of the actual segment locations


type Entry struct {
	TimeStamp uint32
	KeySize   uint32
	ValueSize uint32
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

func (e *Entry) getEntryKVPairSize() uint32 {
	return e.KeySize + e.ValueSize;
}


func (e *Entry) getEntrySize() uint32 {
	return e.getEntryKVPairSize() + 12; // 12 bytes for the headers
}



func (e *Entry) Serialize() []byte {
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

	return bytarr

}


func (e *Entry) Deserialize(bytarr []byte) (error, []byte) {
	offset := 0

	e.TimeStamp = binary.LittleEndian.Uint32(bytarr[offset:])

	offset += 4

	e.KeySize = binary.LittleEndian.Uint32(bytarr[offset:])

	offset += 4

	e.ValueSize = binary.LittleEndian.Uint32(bytarr[offset:])

	if int(e.KeySize)+int(e.ValueSize) != len(bytarr)-12 {
		return utils.ErrInvalidEntry, nil
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

	return nil, bytarr

}