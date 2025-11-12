package core

import (
	"encoding/binary"
	"fmt"

	"github.com/AatirNadim/getMe/server/store/utils"
	"github.com/AatirNadim/getMe/server/utils/logger"
)

// we are dealing with the segment ids instead of the actual segment locations

// the size of an entry instance will be variable, depending on the key and value sizes
type Entry struct {
	TimeStamp int64
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

func CreateEntry(key []byte, value []byte, timeStamp int64) (*Entry, error) {
	logger.Info("Creating entry with key: ", string(key), " and value: ", value)
	return &Entry{
		TimeStamp: timeStamp,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       key,
		Value:     value,
	}, nil
}

func CreateDeletionEntry(key []byte, timeStamp int64) (*Entry, error) {
	return &Entry{
		TimeStamp: timeStamp,
		KeySize:   uint32(len(key)),
		ValueSize: 0,
		Key:       key,
		Value:     nil,
	}, nil
}

func (e *Entry) isDeletionEntry() bool {
	return e.ValueSize == 0
}

func (e *Entry) getEntryKVPairSize() uint32 {
	return e.KeySize + e.ValueSize
}

func getEntryHeaderSize() uint32 {
	return 8 + 4 + 4 // timestamp + keysize + valuesize
}

func (e *Entry) getEntrySize() uint32 {
	return e.KeySize + e.ValueSize + getEntryHeaderSize()
}

func (e *Entry) Serialize() ([]byte, error) {
	logger.Info("Serializing entry with key: ", string(e.Key), " and value size: ", e.ValueSize)

	bytarr := make([]byte, e.getEntrySize())

	offset := 0

	// the format is little endian for all fields since we are targeting x86 compatible architectures

	binary.LittleEndian.PutUint64(bytarr[offset:], uint64(e.TimeStamp))

	offset += 8

	binary.LittleEndian.PutUint32(bytarr[offset:], e.KeySize)

	offset += 4

	binary.LittleEndian.PutUint32(bytarr[offset:], e.ValueSize)

	offset += 4

	copy(bytarr[offset:], e.Key)

	offset += int(e.KeySize)

	if e.ValueSize > 0 {
		copy(bytarr[offset:], e.Value)
	}

	return bytarr, nil

}

func deserializeEntry(bytarr []byte) (*Entry, error) {

	logger.Info("Deserializing entry")

	offset := 0

	e := &Entry{}
	e.TimeStamp = int64(binary.LittleEndian.Uint64(bytarr[offset:]))

	offset += 8

	e.KeySize = binary.LittleEndian.Uint32(bytarr[offset:])

	offset += 4

	e.ValueSize = binary.LittleEndian.Uint32(bytarr[offset:])

	if int(e.KeySize)+int(e.ValueSize) != len(bytarr)-16 {
		fmt.Println("entry.go: 107: Invalid entry: key size + value size does not match the total size. Key size:", e.KeySize, " Value size:", e.ValueSize, " Total size:", len(bytarr))
		return nil, utils.ErrInvalidEntry
	}

	offset += 4

	e.Key = make([]byte, e.KeySize)
	copy(e.Key, bytarr[offset:offset+int(e.KeySize)])

	offset += int(e.KeySize)

	if e.ValueSize > 0 {
		e.Value = make([]byte, e.ValueSize)
		copy(e.Value, bytarr[offset:offset+int(e.ValueSize)])
	} else {
		e.Value = nil
	}

	return e, nil
}

func getEntrySizeFromHeader(header []byte) (uint32, error) {
	if len(header) < int(getEntryHeaderSize()) {
		fmt.Println("entry.go: 130: Invalid header: header size is less than expected size:", len(header), " expected at least:", getEntryHeaderSize())
		return 0, utils.ErrInvalidEntry
	}

	keySize := binary.LittleEndian.Uint32(header[8:12])
	valueSize := binary.LittleEndian.Uint32(header[12:16])

	return getEntryHeaderSize() + keySize + valueSize, nil
}
