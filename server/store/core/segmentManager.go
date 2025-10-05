package core

import (
	"fmt"
	"getMeMod/utils/logger"
	"os"
	"path/filepath"
	"sync"
)

type SegmentManager struct {
	mu         sync.RWMutex
	basePath   string
	segmentMap map[uint32]*Segment
	// stores the index of the next segment to be created
	nextSegmentId uint32
	atomicCounter  *AtomicCounter
}

func NewSegmentManager(basePath string, centralHashTable *HashTable) (*SegmentManager, error) {

	sm := &SegmentManager{
		segmentMap:    make(map[uint32]*Segment),
		basePath:      basePath,
		nextSegmentId: 0,
		// atomicCounter: NewAtomicCounter(0),
	}

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// Load existing segments, if they exist
	if err := sm.populateSegmentMap(basePath, centralHashTable); err != nil {
		return nil, fmt.Errorf("failed to load segments: %w", err)
	}

	// Create active segment if none exists
	if len(sm.segmentMap) == 0 {
		if _, err := sm.createNewSegmentTxn(basePath); err != nil {
			return nil, fmt.Errorf("failed to create active segment: %w", err)
		}
	}

	
	sm.atomicCounter = NewAtomicCounter(sm.nextSegmentId)

	logger.Info("Segment manager atomic counter initialized")
	
	return sm, nil
}

func (sm *SegmentManager) populateSegmentMap(basePath string, centralHashTable *HashTable) error {

	// find all segment files in the base path
	paths, err := filepath.Glob(filepath.Join(basePath, "segment_*.log"))
	if err != nil {
		logger.Error("failed to list segment files basePath", basePath, "error", err)
		return fmt.Errorf("failed to list segment files in %s: %w", basePath, err)
	}
	if paths == nil {
		logger.Warn("no segments found in " + basePath)
		return nil // No segments found is not an error
	}

	logger.Info("Segments already exist, loading them from the disk to the current kv instance...")

	var wg sync.WaitGroup
	ch := make(chan *HashTable, len(paths))

	// for all the paths, open the segment and add it to the segment map, based on their IDs
	for _, path := range paths {
		var id uint32
		_, err := fmt.Sscanf(filepath.Base(path), "segment_%d.log", &id)
		if err != nil {
			return err
		}

		segment, err := OpenSegment(id, basePath)
		if err != nil {
			return err
		}

		logger.Info("opened segment with id:", id, " at path:", path, "segment details:", *segment)

		sm.nextSegmentId = max(sm.nextSegmentId, id)

		// assign the segment mapped to its id
		sm.segmentMap[uint32(id)] = segment

		wg.Add(1)
		// go func(seg *Segment) {
		// 	defer wg.Done()
		// 	ht, err := seg.ReadAllEntries()
		// 	if err != nil {
		// 		logger.Error("Failed to read entries from segment %d: %v", seg.id, err)
		// 		return
		// 	}
		// 	ch <- ht
		// }(segment)

		segment.readAllEntriesAsync(&wg, ch)
		
	}

	// a goroutine to close the channel when all segment reads are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Reducer
	for ht := range ch {
		centralHashTable.Merge(ht)
	}

	// remove all the deletion entries from the central hash table
	// since they are not actual entries
	centralHashTable.DeleteDeletionEntries()

	logger.Success("all the segments have been loaded into the central hash table")

	logger.Info("current segment map : ", sm.segmentMap)

	logger.Info("current hash table : ", centralHashTable.Entries())
	// if the latest entry is a deletion entry, simply remove it from the hash table
	centralHashTable.DeleteDeletionEntries()

	logger.Info("loaded segments from the disk")
	// increment the nextSegmentId to be one more than the max id found
	sm.nextSegmentId += 1
	return nil

}


// create a new segment, add it to the segment map and return it (this operation is not atomic)
func (sm *SegmentManager) createNewSegment(basePath string) (*Segment, error) {
	logger.Info("Creating a new segment with id:", sm.nextSegmentId)

	// create a new segment with id = sm.nextSegmentId
	segment, err := NewSegment(sm.nextSegmentId, basePath)
	if err != nil {
		return nil, err
	}

	// add the new segment to the segment map and increment the nextSegmentId
	sm.segmentMap[sm.nextSegmentId] = segment
	// increment the nextSegmentId for the next segment
	// sm.nextSegmentId += 1

	// use the atomic counter to get the next segment id
	sm.nextSegmentId = sm.atomicCounter.Next()
	return segment, nil
}

// create a new segment, append it to the segment list and return it (this operation is atomic)
func (sm *SegmentManager) createNewSegmentTxn(basePath string) (*Segment, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.createNewSegment(basePath)
}

func (sm *SegmentManager) Append(entry *Entry) (uint32, uint32, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Info("segment manager: Appending entry with key:", string(entry.Key))
	offset, err := sm.appendEntryToLatestSegment(entry)

	if err != nil {
		logger.Error("segment manager: failed to append entry:", err)
		return 0, 0, err
	}

	// Return the segment ID and offset
	return sm.nextSegmentId, offset, nil
}

// reads an entry from a specific segment at a specific offset and returns it along with the offset for the next entry
func (sm *SegmentManager) Read(segmentId uint32, offset uint32) (*Entry, uint32, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	logger.Info("segment manager: Reading entry from segment", segmentId, "at offset", offset)

	if segmentId > sm.nextSegmentId-1 {
		logger.Error("segment manager: segment", segmentId, "does not exist")
		return nil, offset, fmt.Errorf("segment %d does not exist", segmentId)
	}

	segment := sm.segmentMap[segmentId]
	return segment.Get(offset)
}

func (sm *SegmentManager) Update(entry *Entry) (uint32, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Info("segment manager: Updating entry with key:", string(entry.Key))

	offset, err := sm.appendEntryToLatestSegment(entry)

	if err != nil {
		logger.Error("segment manager: failed to update entry:", err)
		return 0, err
	}

	return offset, nil
}

func (sm *SegmentManager) Delete(entry *Entry) (uint32, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Info("segment manager: Deleting entry")

	offset, err := sm.appendEntryToLatestSegment(entry)

	if err != nil {
		return 0, err
	}

	return offset, nil

}

func (sm *SegmentManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for id, segment := range sm.segmentMap {
		segment.file.Close()
		os.Remove(segment.path)
		delete(sm.segmentMap, id)
	}

	sm.nextSegmentId = 0
}

func (sm *SegmentManager) appendEntryToLatestSegment(entry *Entry) (uint32, error) {

	// active id will always hold the id of the next segment to be created
	currentSegment := sm.segmentMap[sm.nextSegmentId-1]

	logger.Info("current segment details: ", *currentSegment)

	if !currentSegment.isSpaceAvailableInCurrentSegment(entry) {
		logger.Info("appendEntryToLatestSegment: No space available in current segment, creating a new segment...")
		newSegment, newSegmentCreationError := sm.createNewSegment(sm.basePath)

		if newSegmentCreationError != nil {
			return 0, newSegmentCreationError
		}

		// Update the current segment to the new segment
		currentSegment = newSegment
	}

	offset, err := currentSegment.Append(entry)

	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (sg *SegmentManager) updateActiveSegmentId(size uint32) (uint32, error) {
	sg.mu.Lock()
	defer sg.mu.Unlock()
	
	currAvailableSegmentId := sg.atomicCounter.Reserve(size)

	sg.nextSegmentId = sg.atomicCounter.Get()

	return currAvailableSegmentId, nil
}

// this will run as a separate goroutine
// perform compaction and returns the compacted hash table to be merged with the central hash table
func (sg *SegmentManager) performCompaction(centralHashTable *HashTable, segments []*Segment, compactedSegmentManager *CompactedSegmentManager) error {


	logger.Debug("Starting compaction for segments, creating a channel")
	resultChan := make(chan *HashTable, len(segments))

	var wg sync.WaitGroup

	// read all the entries in these segmnents and create a new hash table, which contains latest and unique entries in the scope of these segments
	for _, segment := range segments {
		wg.Add(1)
		segment.readAllEntriesAsync(&wg, resultChan)

	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	updatedHashTable := NewHashTable()

	for ht := range resultChan {
		updatedHashTable.Merge(ht)
	}

	logger.Info("updated hash table created")

	// reserve segment ids for the compacted segments
	currAvailableSegmentId, err := sg.updateActiveSegmentId(uint32(len(segments)))

	if err != nil {
		return fmt.Errorf("performCompaction: failed to update active segment id: %w", err)
	}

	logger.Info("active segment id updated to:", currAvailableSegmentId)

	
	

	compactedSegmentManager.currAvailableSegmentId = currAvailableSegmentId
	compactedSegmentManager.maxAvailableSegmentId = currAvailableSegmentId + uint32(len(segments)) - 1
	// we can send the entire segment map, since compaction only deals with inactive segments and it is not being modified
	compactedSegmentManager.originalSegmentMap = sg.segmentMap


	compactedHashTable, err := compactedSegmentManager.populateCompactedSegments(updatedHashTable)


	if err != nil {
		return fmt.Errorf("performCompaction: failed to populate compacted segments: %w", err)
	}


	// bring the segments from the compacted segment manager to the main segment manager
	for id, segment := range compactedSegmentManager.compactedSegmentMap {

		newPath := filepath.Join(sg.basePath, fmt.Sprintf("segment_%d.log", id))
		err := os.Rename(segment.path, newPath)
		if err != nil {
			return fmt.Errorf("performCompaction: failed to move compacted segment file %s to %s: %w", segment.path, newPath, err)
		}
		segment.path = newPath
		
		sg.segmentMap[id] = segment
	}

	// delete the original segments that were compacted
	for _, segment := range segments {
		logger.Info("performCompaction: deleting original segment with id:", segment.id)
		segment.file.Close()
		err := os.Remove(segment.path)
		if err != nil {
			return fmt.Errorf("performCompaction: failed to delete original segment file %s: %w", segment.path, err)
		}
		delete(sg.segmentMap, segment.id)
	}

	logger.Info("performCompaction: deleted original segments")
	
	// merge the compacted hash table into the central hash table 
	centralHashTable.Merge(compactedHashTable)
	
	logger.Debug("Compaction completed. Updated hash table: ", updatedHashTable.Entries())
	return nil

}