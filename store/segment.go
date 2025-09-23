package store

import (
	"encoding/binary"
	"getMeMod/store/utils"
	"os"
	"sync"
)


const DefaultMaxSegmentSize = 1024 * 1024 * 20 // 20 MB

const MaxEntriesPerSegment = 10000


type Segment struct {
	mu sync.RWMutex
	id int
	path string
	file *os.File
	entryCount int
	size int
	isActive bool
	maxCount int
	maxSize int
	
}


func (segment *Segment) Append(entry* Entry) (uint32, error) {
	segment.mu.Lock()

	defer segment.mu.Unlock()

	if(segment.size > segment.maxSize || segment.entryCount >= segment.maxCount) {
		return 0, utils.ErrSegmentFull
	}


	data := entry.Serialize()

	segmentOffset := uint32(segment.size)

	_, err := segment.file.Write(data) // append bytes to the file
	if err != nil {
		return 0, err
	}

	segment.size += len(data)
	segment.entryCount += 1

	return segmentOffset, nil
}


func (segment *Segment) Get(pos uint32) (error, *Entry) {
	segment.mu.RLock()
	defer segment.mu.RUnlock()

	if pos >= uint32(segment.size) {
		return utils.ErrInvalidEntry, nil
	}

	// read the entry header first to determine sizes
	header := make([]byte, 12)
	_, err := segment.file.ReadAt(header, int64(pos))
	if err != nil {
		return err, nil
	}

	keySize := binary.LittleEndian.Uint32(header[4:8])
	valueSize := binary.LittleEndian.Uint32(header[8:12])

	serializedEntry := make([]byte, 12+keySize+valueSize)

	entry, err := DeserializeEntry(serializedEntry)
	if err != nil {
		return err, nil
	}

	return nil, entry
}


