package core

import (
	"fmt"
	"getMeMod/server/store/utils/constants"
	"os"
	"path/filepath"
	"sync"
)


type CompactedSegmentManager struct {
	mu         sync.RWMutex
	basePath   string
	currAvailableSegmentId uint32
	maxAvailableSegmentId uint32

	// to keep track of the original segments that are to be compacted, we can fetch the active segment ids from this
	activeSegment *Segment
	compactedHashTable *HashTable
	originalSegmentMap map[uint32]*Segment
	compactedSegmentMap map[uint32]*Segment
}


func NewCompactedSegmentManager(basePath string, currAvailableSegmentId, maxAvailableSegmentId uint32, compactedHashTable *HashTable) (*CompactedSegmentManager, error) {

	csm := &CompactedSegmentManager{
		basePath:      basePath,
		currAvailableSegmentId: currAvailableSegmentId,
		maxAvailableSegmentId: maxAvailableSegmentId,
		compactedHashTable: compactedHashTable,
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}


	// TODO: check whether this is required
	// Load existing compacted segments, if they exist
	if _, err := csm.populateCompactedSegments(csm.compactedHashTable); err != nil {
		return nil, fmt.Errorf("failed to load compacted segments: %w", err)
	}

	return csm, nil
}



func (csm *CompactedSegmentManager) populateCompactedSegments(compactedHashTable *HashTable) (*HashTable, error) {
	
	for key, hashTableEntry := range compactedHashTable.Entries() {

		entry, _, err := csm.originalSegmentMap[hashTableEntry.SegmentId].Get(hashTableEntry.Offset)

		if err != nil {
			return nil, fmt.Errorf("failed to get entry from original segment: %w", err)
		}
		
		offset, err := csm.appendEntryToActiveCompactedSegment(entry)

		updatedhashTableEntry := hashTableEntry;
		updatedhashTableEntry.SegmentId = csm.currAvailableSegmentId
		updatedhashTableEntry.Offset = offset
		
		compactedHashTable.PutEntry(key, updatedhashTableEntry)
	}

	return compactedHashTable, nil
}


// this is supposed to return the segment id and the offset of the appended entry
func (csm *CompactedSegmentManager) appendEntryToActiveCompactedSegment(entry *Entry) (uint32, error) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	if csm.activeSegment == nil || !csm.activeSegment.isSpaceAvailableInCurrentSegment(entry) {
		if err := csm.createdNewSegment(); err != nil {
			return 0, err
		}
	}

	// we now have the latest active segment to append entry
	offset, err := csm.activeSegment.Append(entry)

	if err != nil {
		return 0, err
	}

	return offset, nil
}

func (csm *CompactedSegmentManager) createdNewSegment() error {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	path := filepath.Join(csm.basePath, fmt.Sprintf("segment_%d.log", csm.currAvailableSegmentId))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return err
	}

	csm.activeSegment = &Segment{
		// for the id, use the currently available segment id
		id:       csm.currAvailableSegmentId,
		path:     path,
		file:     file,
		isActive: true,
		maxCount: constants.MaxEntriesPerSegment,
		maxSize:  constants.DefaultMaxSegmentSize,
	}
	
	csm.compactedSegmentMap[csm.currAvailableSegmentId] = csm.activeSegment

	csm.currAvailableSegmentId += 1
	return nil
}