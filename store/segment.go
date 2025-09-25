package store

import (
	"encoding/binary"
	"fmt"
	"getMeMod/store/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
)


const DefaultMaxSegmentSize = 1024 * 1024 * 20 // 20 MB

const MaxEntriesPerSegment = 10000


// represents a log segment file, stored on the disk
type Segment struct {
	mu sync.RWMutex
	id uint32
	path string
	file *os.File
	entryCount int
	size int
	isActive bool
	maxCount int
	maxSize int
	
}

// takes in the id of the new segment to be created, and the base path where it should be created
func NewSegment(id uint32, basePath string) (*Segment, error) {
	path := filepath.Join(basePath, fmt.Sprintf("segment_%d.log", id))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}
	
	segment := &Segment{
		id:         id,
		path:      path,
		file:      file,
		isActive:  true,
		maxCount:  MaxEntriesPerSegment,
		maxSize:   DefaultMaxSegmentSize,
		
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
		id:         id,
		path:      path,
		file:      file,
		isActive:  true,
		maxCount:  MaxEntriesPerSegment,
		maxSize:   DefaultMaxSegmentSize,

	}

	return segment, nil
}


// appends the entry to the segment file, returns the offset at which the entry was added
func (segment *Segment) Append(entry *Entry) (uint32, error) {
	segment.mu.Lock()

	defer segment.mu.Unlock()


	// no space left in the current segment to add new entries
	if(segment.size > segment.maxSize || segment.entryCount >= segment.maxCount) {
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

	segment.size += len(data)
	segment.entryCount += 1

	// return the starting position of the newly added entry in the segment file
	return segmentOffset, nil
}


// takes in the starting position of the entry in the segment file and returns the entry
func (segment *Segment) Get(pos uint32) (*Entry, error) {
	segment.mu.RLock()
	defer segment.mu.RUnlock()

	// check if the position is valid
	if pos >= uint32(segment.size) {
		return nil, utils.ErrInvalidEntry
	}

	// read the entry header first to determine sizes
	header := make([]byte, 12)
	_, err := segment.file.ReadAt(header, int64(pos))
	if err != nil {
		return nil, err
	}

	// extract key and value sizes from the header
	keySize := binary.LittleEndian.Uint32(header[4:8])

	// extract value size from the header
	valueSize := binary.LittleEndian.Uint32(header[8:12])

	serializedEntry := make([]byte, 12+keySize+valueSize)

	entry, err := DeserializeEntry(serializedEntry)
	if err != nil {
		return nil, err
	}

	// return the entry in the given segment at the given offset
	return entry, nil
}


func isSpaceAvailableInCurrentSegment(segment *Segment, entry *Entry) bool {

	log.Printf("Current segment size: %d, max size: %d, entry count: %d, max count: %d, new entry size: %d\n", segment.size, segment.maxSize, segment.entryCount, segment.maxCount, entry.getEntrySize())
	return segment.size + int(entry.getEntrySize()) <= segment.maxSize && segment.entryCount < segment.maxCount
}

// creates a deletion entry for the given key
func (segment *Segment) CreateDeletionEntry(key []byte) (*Entry, error) {

	// create a deletion entry for the given key
	log.Println("segment file: Creating deletion entry for key:", string(key))
	entry, err := CreateDeletionEntry(key)
	if err != nil {
		return nil, err
	}
	return entry, nil
}