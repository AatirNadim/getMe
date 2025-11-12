package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AatirNadim/getMe/server/store/utils/constants"
	"github.com/AatirNadim/getMe/server/utils/logger"
)

type CompactionResult struct {
	CompactedHashTable *HashTable
	OldSegmentIds      []uint32
}

type CompactedSegmentManager struct {
	basePath               string
	nextAvailableSegmentId uint32
	maxAvailableSegmentId  uint32
	// to keep track of the original segments that are to be compacted, we can fetch the active segment ids from this
	activeSegment       *Segment
	originalSegmentMap  map[uint32]*Segment
	compactedSegmentMap map[uint32]*Segment
}

func NewCompactedSegmentManager(basePath string) (*CompactedSegmentManager, error) {

	csm := &CompactedSegmentManager{
		basePath:            basePath,
		compactedSegmentMap: make(map[uint32]*Segment),
	}

	// remove existing directory if it exists
	err := os.RemoveAll(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to remove dir: %w", err)
	}

	// create base directory
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// // TODO: check whether this is required
	// // Load existing compacted segments, if they exist
	// if _, err := csm.populateCompactedSegments(csm.compactedHashTable); err != nil {
	// 	return nil, fmt.Errorf("failed to load compacted segments: %w", err)
	// }

	return csm, nil
}

func (csm *CompactedSegmentManager) populateCompactedSegments(compactedHashTable *HashTable) (*HashTable, error) {

	logger.Info("Populating compacted segments from base path:", csm.basePath)

	for key, hashTableEntry := range compactedHashTable.Entries() {

		entry, _, err := csm.originalSegmentMap[hashTableEntry.SegmentId].Get(hashTableEntry.Offset)

		if err != nil {
			return nil, fmt.Errorf("failed to get entry from original segment: %w", err)
		}

		offset, err := csm.appendEntryToActiveCompactedSegment(entry)

		if err != nil {
			return nil, fmt.Errorf("failed to append entry to active compacted segment: %w", err)
		}

		updatedhashTableEntry := hashTableEntry
		updatedhashTableEntry.SegmentId = csm.nextAvailableSegmentId - 1 // the segment id of the active compacted segment
		updatedhashTableEntry.Offset = offset

		compactedHashTable.PutEntry(key, updatedhashTableEntry)
	}

	return compactedHashTable, nil
}

// this is supposed to return the segment id and the offset of the appended entry
func (csm *CompactedSegmentManager) appendEntryToActiveCompactedSegment(entry *Entry) (uint32, error) {

	if csm.activeSegment == nil || !csm.activeSegment.isSpaceAvailableInCurrentSegment(entry) {
		if err := csm.createNewSegment(); err != nil {
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

func (csm *CompactedSegmentManager) createNewSegment() error {
	path := filepath.Join(csm.basePath, fmt.Sprintf("segment_%d.log", csm.nextAvailableSegmentId))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		return err
	}

	logger.Info("Created new compacted segment with id:", csm.nextAvailableSegmentId, " at path:", path)

	csm.activeSegment = &Segment{
		// for the id, use the currently available segment id
		id:       csm.nextAvailableSegmentId,
		path:     path,
		file:     file,
		isActive: true,
		maxCount: constants.MaxEntriesPerSegment,
		maxSize:  constants.DefaultMaxSegmentSize,
	}

	csm.compactedSegmentMap[csm.nextAvailableSegmentId] = csm.activeSegment

	csm.nextAvailableSegmentId += 1

	return nil
}

func (csm *CompactedSegmentManager) clearManager() {

	csm.compactedSegmentMap = make(map[uint32]*Segment)
	csm.nextAvailableSegmentId = 0
	csm.maxAvailableSegmentId = 0
	csm.originalSegmentMap = make(map[uint32]*Segment)
	csm.activeSegment = nil
}
