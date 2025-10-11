package core

import (
	"fmt"
	"getMeMod/server/store/utils"
	"getMeMod/server/store/utils/constants"
	"getMeMod/utils/logger"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// represents a log segment file, stored on the disk
type Segment struct {
	// mu         sync.RWMutex
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
		maxCount: constants.MaxEntriesPerSegment,
		maxSize:  constants.DefaultMaxSegmentSize,
	}

	return segment, nil
}

func OpenSegment(id uint32, basePath string) (*Segment, error) {

	// Construct the file path for the segment

	logger.Info("Opening segment with id:", id, "at the base path:", basePath)

	path := filepath.Join(basePath, fmt.Sprintf("segment_%d.log", id))

	// Open the segment file in read-write and append mode, returns the pointer to the file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return nil, utils.ErrFileNotFoundOrNotAccessible
	}

	fileInfo, statErr := file.Stat()

	if statErr != nil {
		return nil, fmt.Errorf("failed to get file info for segment %s: %w", path, statErr)
	}

	segment := &Segment{
		id:       id,
		path:     path,
		file:     file,
		isActive: true,
		maxCount: constants.MaxEntriesPerSegment,
		maxSize:  constants.DefaultMaxSegmentSize,
		size:     uint32(fileInfo.Size()),
	}

	return segment, nil
}

// appends the entry to the segment file, returns the offset at which the entry was added
func (segment *Segment) Append(entry *Entry) (uint32, error) {
	// segment.mu.Lock()

	// defer segment.mu.Unlock()

	// no space left in the current segment to add new entries
	if segment.size > segment.maxSize || segment.entryCount >= segment.maxCount {
		return 0, utils.ErrSegmentFull
	}

	data, err := entry.Serialize()

	logger.Info("Serialized data: ", data)

	if err != nil {
		return 0, err
	}

	// calculate the end of the file to include the entry
	segmentOffset := uint32(segment.size)

	_, writeError := segment.file.Write(data) // append bytes to the file
	if writeError != nil {
		return 0, writeError
	}

	// updating the segment size here
	segment.size += uint32(len(data))
	segment.entryCount += 1

	logger.Info("segment.go : current segment size: ", len(data), segment.size, ", current entry count: ", segment.entryCount)

	// return the starting position of the newly added entry in the segment file
	return segmentOffset, nil
}

func (segment *Segment) AppendBuffer(buffer []byte) error {
	// segment.mu.Lock()
	// defer segment.mu.Unlock()

	// for now, we can append the buffer directly to the current segment, irrespetive of its size

	// if uint32(len(buffer)) > (segment.maxSize - segment.size) {
	// 	return utils.ErrSegmentFull
	// }

	_, writeError := segment.file.Write(buffer)
	if writeError != nil {
		return writeError
	}

	segment.size += uint32(len(buffer))
	// Note: entryCount is not updated here as we are writing a raw buffer.
	// The caller (BatchPut) is responsible for managing entry semantics.

	return nil
}

// takes in the starting position of the entry in the segment file and returns the entry and the offset for the next entry
func (segment *Segment) Get(offset uint32) (*Entry, uint32, error) {
	// segment.mu.RLock()
	// defer segment.mu.RUnlock()

	// check if the position is valid
	if offset >= uint32(segment.size) {
		return nil, offset, utils.ErrInvalidEntry
	}

	serializedEntry, newOffset, err := segment.getSerializedEntryFromOffset(offset)

	if err != nil {
		return nil, offset, err
	}

	entry, err := deserializeEntry(serializedEntry)
	if err != nil {
		return nil, offset, err
	}

	// return the entry in the given segment at the given offset and the offset for the next entry
	return entry, newOffset, nil
}

// this function checks whether there is any space available in the current segment to add a new entry
// this is useful when we want to append an entry with size greater than the max segment size
// in the alternative, new segment creation will be cascaded
func (segment *Segment) isSpaceAvailableInCurrentSegment(entry *Entry) bool {

	logger.Info("is space available: Current segment size: %d, max size: %d, entry count: %d, max count: %d, new entry size: %d\n", segment.size, segment.maxSize, segment.entryCount, segment.maxCount, entry.getEntrySize())
	return segment.size <= segment.maxSize && segment.entryCount < segment.maxCount
}

// reads all entries from the segment file and returns a hashtable with the key and its corresponding segment id, offset and timestamp, to be used for map-reduce operations
func (segment *Segment) ReadAllEntries() (*HashTable, error) {
	// segment.mu.RLock()
	// defer segment.mu.RUnlock()

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

		entry, desErr := deserializeEntry(serializedEntry)
		if desErr != nil {
			// Log this, but continue if possible, as it might be a partial write
			logger.Error("Failed to deserialize entry at offset %d: %v", offset, desErr)
			offset = nextOffset
			continue
		}

		entryKey := convertBytesToString(entry.Key)

		// at the time of creation of new hashtable, deletion entries remove the key if present
		// if entry.isDeletionEntry() {
		// 	ht.Delete(entryKey)
		// } else {

		// for now we will keep the deletion entries in the hash table as well to keep track of the latest entries
		ht.Put(entryKey, segment.id, offset, entry.TimeStamp, entry.ValueSize)
		// }

		// updating the offset to point to the next entry
		offset = nextOffset
	}

	return ht, nil
}

func (segment *Segment) readAllEntriesAsync(wg *sync.WaitGroup, resultChan chan<- *HashTable) {
	defer wg.Done()
	ht, err := segment.ReadAllEntries()
	if err != nil {
		logger.Error("Failed to read entries from segment %d: %v", segment.id, err)
		return
	}
	resultChan <- ht
}

// reads an entry from a specific offset in the segment file and returns the serialized bytes of the entry along with the offset for the next entry
func (sg *Segment) getSerializedEntryFromOffset(offset uint32) ([]byte, uint32, error) {
	// sg.mu.RLock()
	// defer sg.mu.RUnlock()

	if offset >= uint32(sg.size) {
		return nil, offset, io.EOF
	}

	// maxSize:  DefaultMaxSegmentSize,
	header := make([]byte, getEntryHeaderSize())
	_, err := sg.file.ReadAt(header, int64(offset))
	if err != nil {
		if err == io.EOF {
			return nil, offset, io.EOF // Reached the end of the file
		}
		return nil, offset, fmt.Errorf("error reading entry header at offset %d: %w", offset, err)
	}
	entrySize, err := getEntrySizeFromHeader(header)

	if err != nil {
		return nil, offset, fmt.Errorf("error getting entry size from header at offset %d: %w", offset, err)
	}

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
