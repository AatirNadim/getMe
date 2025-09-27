package store

import (
	"encoding/binary"
	"fmt"
	"getMeMod/store/logger"
	"getMeMod/store/utils"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const DefaultMaxSegmentSize = 1024 * 1024 * 20 // 20 MB

const MaxEntriesPerSegment = 10000

// represents a log segment file, stored on the disk
type Segment struct {
	mu         sync.RWMutex
	id         uint32
	path       string
	file       *os.File
	entryCount uint32
	size       uint32
	isActive   bool
	maxCount   uint32
	maxSize    uint32
}

// takes in the id of the new segment to be created, and the base path where it should be created
func NewSegment(id uint32, basePath string) (*Segment, error) {
	path := filepath.Join(basePath, fmt.Sprintf("segment_%d.log", id))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	segment := &Segment{
		id:       id,
		path:     path,
		file:     file,
		isActive: true,
		maxCount: MaxEntriesPerSegment,
		maxSize:  DefaultMaxSegmentSize,
	}

	return segment, nil
}

func OpenSegment(id uint32, basePath string) (*Segment, error) {

	// Construct the file path for the segment
	path := filepath.Join(basePath, fmt.Sprintf("segment_%d.log", id))

	// Open the segment file in read-write and append mode, returns the pointer to the file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return nil, utils.ErrFileNotFoundOrNotAccessible
	}

	segment := &Segment{
		id:       id,
		path:     path,
		file:     file,
		isActive: true,
		maxCount: MaxEntriesPerSegment,
		maxSize:  DefaultMaxSegmentSize,
	}

	

	return segment, nil
}

// appends the entry to the segment file, returns the offset at which the entry was added
func (segment *Segment) Append(entry *Entry) (uint32, error) {
	segment.mu.Lock()

	defer segment.mu.Unlock()

	// no space left in the current segment to add new entries
	if segment.size > segment.maxSize || segment.entryCount >= segment.maxCount {
		return 0, utils.ErrSegmentFull
	}

	data, err := entry.Serialize()
	if err != nil {
		return 0, err
	}

	// calculate the end of the file to include the entry
	segmentOffset := uint32(segment.size)

	_, writeError := segment.file.Write(data) // append bytes to the file
	if writeError != nil {
		return 0, writeError
	}

	segment.size += uint32(len(data))
	segment.entryCount += 1

	// return the starting position of the newly added entry in the segment file
	return segmentOffset, nil
}

// takes in the starting position of the entry in the segment file and returns the entry and the offset for the next entry
func (segment *Segment) Get(offset uint32) (*Entry, uint32, error) {
	segment.mu.RLock()
	defer segment.mu.RUnlock()

	// check if the position is valid
	if offset >= uint32(segment.size) {
		return nil, offset, utils.ErrInvalidEntry
	}

	// read the entry header first to determine sizes
	header := make([]byte, 12)
	_, err := segment.file.ReadAt(header, int64(offset))
	if err != nil {
		return nil, offset, err
	}

	serializedEntry, newOffset, err := segment.getSerializedEntryFromOffset(offset)

	if err != nil {
		return nil, offset, err
	}

	entry, err := DeserializeEntry(serializedEntry)
	if err != nil {
		return nil, offset, err
	}

	// return the entry in the given segment at the given offset and the offset for the next entry
	return entry, newOffset, nil
}

func isSpaceAvailableInCurrentSegment(segment *Segment, entry *Entry) bool {

	logger.Info("Current segment size: %d, max size: %d, entry count: %d, max count: %d, new entry size: %d\n", segment.size, segment.maxSize, segment.entryCount, segment.maxCount, entry.getEntrySize())
	return segment.size+entry.getEntrySize() <= segment.maxSize && segment.entryCount < segment.maxCount
}

// creates a deletion entry for the given key
func (segment *Segment) CreateDeletionEntry(key []byte) (*Entry, error) {

	// create a deletion entry for the given key
	logger.Info("segment file: Creating deletion entry for key:", string(key))
	entry, err := CreateDeletionEntry(key)
	if err != nil {
		return nil, err
	}
	return entry, nil
}
// reads all entries from the segment file and returns a hashtable with the key and its corresponding segment id, offset and timestamp, to be used for map-reduce operations
func (segment *Segment) ReadAllEntries() (*HashTable, error) {
	segment.mu.RLock()
	defer segment.mu.RUnlock()

	ht := NewHashTable()
	offset := uint32(0)

	// run this loop till we reach the end of the file
	for {
		serializedEntry, nextOffset, err := segment.getSerializedEntryFromOffset(offset)

		if err != nil {
			if err == io.EOF {
				break // Reached end of segment file
			}
			return nil, err
		}

		entry, desErr := DeserializeEntry(serializedEntry)
		if desErr != nil {
			// Log this, but continue if possible, as it might be a partial write
			logger.Error("Failed to deserialize entry at offset %d: %v", offset, desErr)
			offset = nextOffset
			continue
		}

		entryKey := convertBytesToString(entry.Key)

		if entry.IsDeletionEntry() {
			ht.Delete(entryKey)
		} else {
			ht.Put(entryKey, segment.id, offset, entry.TimeStamp)
		}

		// updating the offset to point to the next entry
		offset = nextOffset
	}

	return ht, nil
}

// reads an entry from a specific offset in the segment file and returns the serialized bytes of the entry along with the offset for the next entry
func (sg *Segment) getSerializedEntryFromOffset(offset uint32) ([]byte, uint32, error) {
	sg.mu.RLock()
	defer sg.mu.RUnlock()

	if offset >= uint32(sg.size) {
		return nil, offset, io.EOF
	}

	// maxSize:  DefaultMaxSegmentSize,
	header := make([]byte, 12)
	_, err := sg.file.ReadAt(header, int64(offset))
	if err != nil {
		if err == io.EOF {
			return nil, offset, io.EOF // Reached the end of the file
		}
		return nil, offset, fmt.Errorf("error reading entry header at offset %d: %w", offset, err)
	}

	keySize := binary.LittleEndian.Uint32(header[4:8])
	valueSize := binary.LittleEndian.Uint32(header[8:12])
	entrySize := 12 + keySize + valueSize

	serializedEntry := make([]byte, entrySize)
	_, err = sg.file.ReadAt(serializedEntry, int64(offset))
	if err != nil {
		if err == io.EOF {
			return nil, offset, io.EOF // Reached the end of the file
		}
		return nil, offset, fmt.Errorf("error reading full entry at offset %d: %w", offset, err)
	}
	offset = offset + uint32(entrySize)
	return serializedEntry, offset, nil
}

func convertStringToBytes(str string) []byte {
	return []byte(str)
}

func convertBytesToString(b []byte) string {
	return string(b)
}