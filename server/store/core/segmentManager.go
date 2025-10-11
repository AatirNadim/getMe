package core

import (
	"fmt"
	"getMeMod/server/store/utils/constants"
	"getMeMod/utils/logger"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type SegmentManager struct {
	mu         sync.RWMutex
	basePath   string
	segmentMap map[uint32]*Segment
	// stores the index of the next segment to be created
	// nextSegmentId uint32
	// atomicCounter  *AtomicCounter
	nextSegmentCounter      *AtomicCounter
	compactedSegmentManager *CompactedSegmentManager
	isCompacting            atomic.Bool
	compactionResultChannel chan *CompactionResult
}

func NewSegmentManager(basePath, compactedBasePath string, centralHashTable *HashTable, compactionResultChannel chan *CompactionResult) (*SegmentManager, error) {

	compactedSegmentManager, err := NewCompactedSegmentManager(compactedBasePath)
	if err != nil {
		panic(err)
	}

	sm := &SegmentManager{
		segmentMap: make(map[uint32]*Segment),
		basePath:   basePath,
		// nextSegmentId: 0,
		nextSegmentCounter:      NewAtomicCounter(0),
		compactedSegmentManager: compactedSegmentManager,
		compactionResultChannel: compactionResultChannel,
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

	maxSegmentId := uint32(0)

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

		maxSegmentId = max(maxSegmentId, id)

		// assign the segment mapped to its id
		sm.segmentMap[uint32(id)] = segment

		wg.Add(1)

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
	sm.nextSegmentCounter.Set(maxSegmentId + 1)
	return nil

}

// gets the total no of segments (this operation is not atomic)
func (sg *SegmentManager) TotalSegments() uint32 {
	// sg.mu.RLock()
	// defer sg.mu.RUnlock()
	logger.Debug("totalSegments")

	return uint32(len(sg.segmentMap))
}

// with the given segment id, create a new segment, add it to the segment map and return it (this operation is not atomic)
func (sm *SegmentManager) createNewSegment(basePath string) (*Segment, error) {

	segmentId := sm.nextSegmentCounter.Get()

	logger.Info("Creating a new segment with id:", segmentId)

	// create a new segment with id = sm.nextSegmentId
	segment, err := NewSegment(segmentId, basePath)
	if err != nil {
		return nil, err
	}

	// add the new segment to the segment map and increment the nextSegmentId
	sm.segmentMap[segmentId] = segment
	// increment the nextSegmentId for the next segment
	// sm.nextSegmentId += 1

	sm.nextSegmentCounter.Next()

	// use the atomic counter to get the next segment id
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
	return sm.nextSegmentCounter.Get() - 1, offset, nil
}

// reads an entry from a specific segment at a specific offset and returns it along with the offset for the next entry
func (sm *SegmentManager) Read(segmentId uint32, offset uint32) (*Entry, uint32, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	logger.Info("segment manager: Reading entry from segment", segmentId, "at offset", offset)

	if segmentId > sm.nextSegmentCounter.Get()-1 {
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

// FlushBuffer writes a buffer of serialized entries to the active segment.
// It returns a map of keys to their new disk locations (SegmentID and Offset).
// this is not an locked operation, i.e. no mutex is held during the operation
func (sm *SegmentManager) FlushBuffer(buffer []byte, entries []*Entry) ([]*FlushResult, error) {
	// sm.mu.Lock()
	// defer sm.mu.Unlock()

	logger.Info("segment manager: Flushing buffer with size:", len(buffer))

	// nextsegmentcounter is referenced only once in the function body, since it is atomic
	currentSegment := sm.segmentMap[sm.nextSegmentCounter.Get()-1]

	// Check if the buffer can fit in the current segment.
	// The store is now responsible for chunking, but we retain this check as a safeguard.
	if uint32(len(buffer)) > (constants.DefaultMaxSegmentSize - currentSegment.size) {
		// This case should ideally not be hit if the store chunks correctly.
		// We will create a new segment to handle this oversized buffer.
		logger.Info("FlushBuffer: Buffer too large for current segment, creating a new one.")
		newSegment, err := sm.createNewSegment(sm.basePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create new segment for flush: %w", err)
		}
		currentSegment = newSegment
	}

	startOffset := currentSegment.size
	if err := currentSegment.AppendBuffer(buffer); err != nil {
		return nil, fmt.Errorf("failed to write buffer to segment: %w", err)
	}

	// flushResults := make(map[string]*FlushResult)


	flushResults := make([]*FlushResult, 0, len(entries))
	currentOffset := int64(startOffset)
	for _, entry := range entries {
		flushResults = append(flushResults, &FlushResult{
			SegmentID: int(currentSegment.id),
			Offset:    currentOffset,
		})
		currentOffset += int64(entry.getEntrySize())
	}

	logger.Debug("FlushBuffer: Flush results:", flushResults)

	return flushResults, nil
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

func (sm *SegmentManager) Clear() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Info("Clearing all segments from the segment manager")

	for id, segment := range sm.segmentMap {
		segment.file.Close()
		os.Remove(segment.path)
		delete(sm.segmentMap, id)
	}

	sm.nextSegmentCounter.Set(0)

	if _, err := sm.createNewSegment(sm.basePath); err != nil {
		return fmt.Errorf("failed to create active segment: %w", err)
	}
	return nil
}

func (sm *SegmentManager) appendEntryToLatestSegment(entry *Entry) (uint32, error) {

	// active id will always hold the id of the next segment to be created

	logger.Info("appendEntryToLatestSegment: appending entry, current nextsegment counter --> ", sm.nextSegmentCounter.Get())
	currentSegment := sm.segmentMap[sm.nextSegmentCounter.Get()-1]

	logger.Debug("current active segment details: ", *currentSegment)

	// logger.Info("current segment details: ", *currentSegment)

	if !currentSegment.isSpaceAvailableInCurrentSegment(entry) {
		logger.Info("appendEntryToLatestSegment: No space available in current segment, creating a new segment...")

		// a new segment is supposed to be added here
		// here, we will trigger compaction on inactive segments as well

		totalSegments := sm.TotalSegments()

		if totalSegments <= constants.ThresholdForCompaction {
			logger.Info("Total segments do not exceed the threshold for compaction, skipping compaction.")

		} else {
			// reserve segment ids for the compacted segments
			segmentsForCompaction := sm.getSegmentsForCompaction()
			logger.Debug("segments selected for compaction: ", segmentsForCompaction)
			currAvailableSegmentId, err := sm.updateActiveSegmentId(uint32(len(segmentsForCompaction)))

			if err != nil {
				logger.Error("appendEntryToLatestSegment: failed to update active segment id: ", err)
				return 0, err
			}

			go sm.PerformCompaction(segmentsForCompaction, currAvailableSegmentId)

		}

		newSegment, newSegmentCreationError := sm.createNewSegment(sm.basePath)

		if newSegmentCreationError != nil {
			return 0, newSegmentCreationError
		}

		// Update the current segment to the new segment
		currentSegment = newSegment
		// update the nextsegment counter to the next id to be used for segment creation

	}

	offset, err := currentSegment.Append(entry)

	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (sg *SegmentManager) updateActiveSegmentId(size uint32) (uint32, error) {
	// sg.mu.Lock()
	// defer sg.mu.Unlock()

	currAvailableSegmentId := sg.nextSegmentCounter.Reserve(size)

	// sg.nextSegmentId = sg.atomicCounter.Get()

	return currAvailableSegmentId, nil
}

func (sg *SegmentManager) getSegmentsForCompaction() []*Segment {
	// sg.mu.RLock()
	// defer sg.mu.RUnlock()

	var segments []*Segment

	activeSegmentId := sg.nextSegmentCounter.Get() - 1

	for segmentId, segment := range sg.segmentMap {
		if len(segments) >= constants.TotalSegmentsToCompactAtOnce {
			break
		}
		if segmentId != activeSegmentId {
			segments = append(segments, segment)
		}
	}

	return segments
}

func (sm *SegmentManager) getScopedSegmentMap(segments []*Segment) map[uint32]*Segment {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	scopedSegmentMap := make(map[uint32]*Segment)
	for _, segment := range segments {
		scopedSegmentMap[segment.id] = segment
	}
	return scopedSegmentMap
}

// this will run as a separate goroutine
// perform compaction and returns the compacted hash table to be merged with the central hash table
func (sm *SegmentManager) PerformCompaction(segments []*Segment, currAvailableSegmentId uint32) {

	if !sm.isCompacting.CompareAndSwap(false, true) {
		logger.Info("Compaction is already in progress, skipping this trigger.")
		return
	}
	defer sm.isCompacting.Store(false)

	// routineFile, err := os.OpenFile("compaction_routine_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// if err != nil {
	// 	logger.Error("Failed to open routine log file:", err)
	// 	return
	// }
	// defer routineFile.Close()

	logger.Info("initiating compaction process")
	// fmt.Fprintf(routineFile, "initiating compaction process, timestamp: %v\n", time.Now())

	// logger.Info("Total segments in the segment manager:", totalSegments)

	// segments := sg.getSegmentsForCompaction()

	// logFile, err := os.OpenFile("compaction_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	logger.Error("Failed to open log file:", err)
	// 	return
	// }
	// defer logFile.Close()

	if len(segments) == 0 {
		logger.Info("No segments available for compaction")
		return
	}

	logger.Debug("Starting compaction for segments, creating a channel")

	// fmt.Fprintf(logFile, "Starting compaction for segments, creating a channel\n")
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

	// if err != nil {
	// 	logger.Error("performCompaction: failed to update active segment id: ", err)
	// 	return
	// }

	// fmt.Fprintf(logFile, "active segment id updated to: %d\n", currAvailableSegmentId)
	logger.Info("active segment id updated to:", currAvailableSegmentId)

	sm.compactedSegmentManager.nextAvailableSegmentId = currAvailableSegmentId
	sm.compactedSegmentManager.maxAvailableSegmentId = currAvailableSegmentId + uint32(len(segments)) - 1
	// we can send the entire segment map, since compaction only deals with inactive segments and it is not being modified

	sm.compactedSegmentManager.originalSegmentMap = sm.getScopedSegmentMap(segments)

	compactedHashTable, err := sm.compactedSegmentManager.populateCompactedSegments(updatedHashTable)

	if err != nil {
		logger.Error("performCompaction: failed to populate compacted segments: ", err)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// fmt.Fprintf(logFile, "\n==================\n")

	// for id, segment := range sm.compactedSegmentManager.compactedSegmentMap {
	// 	fmt.Fprintf(logFile, "\n\nCompacted segment created with id: %d at path: %s\n\n", id, segment.path)

	// }
	// fmt.Fprintf(logFile, "\n==================\n")
	// bring the segments from the compacted segment manager to the main segment manager
	// this needs to be an atomic operation
	for id, segment := range sm.compactedSegmentManager.compactedSegmentMap {

		newPath := filepath.Join(sm.basePath, fmt.Sprintf("segment_%d.log", id))
		err := os.Rename(segment.path, newPath)
		if err != nil {
			logger.Error(fmt.Errorf("performCompaction: failed to move compacted segment file %s to %s: %w", segment.path, newPath, err))
			return
		}
		segment.path = newPath

		sm.segmentMap[id] = segment
	}

	logger.Debug("clearing the compacted segment map")
	sm.compactedSegmentManager.clearManager()

	sm.compactionResultChannel <- &CompactionResult{
		CompactedHashTable: compactedHashTable,
		OldSegmentIds: func() []uint32 {
			var ids []uint32
			for _, segment := range segments {
				ids = append(ids, segment.id)
			}
			return ids
		}(),
	}

	// merge the compacted hash table into the central hash table

	// fmt.Fprintf(logFile, "Merging compacted hash table into central hash table\n")
	// centralHashTable.Merge(compactedHashTable)

	// centralHashTable.DeleteDeletionEntries()

	// delete the original segments that were compacted
	// for _, segment := range segments {
	// 	logger.Info("performCompaction: deleting original segment with id:", segment.id)
	// 	fmt.Fprintf(logFile, "performCompaction: deleting original segment with id: %d\n", segment.id)
	// 	segment.file.Close()
	// 	err := os.Remove(segment.path)
	// 	if err != nil {
	// 		logger.Error(fmt.Errorf("performCompaction: failed to delete original segment file %s: %w", segment.path, err))
	// 		return
	// 	}
	// 	delete(sg.segmentMap, segment.id)
	// }

	logger.Info("performCompaction: deleted original segments")

	logger.Debug("Compaction completed. Updated hash table: ", compactedHashTable.Entries())
	// fmt.Fprintf(logFile, "Compaction completed. Updated hash table: %v\n", compactedHashTable.Entries())
	// fmt.Fprintf(routineFile, "compaction process completed, timestamp: %v\n", time.Now())

}

func (sm *SegmentManager) DeleteOldSegments(oldSegmentIds []uint32) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, segmentId := range oldSegmentIds {
		if segment, exists := sm.segmentMap[segmentId]; exists {
			logger.Info("Deleting old segment with id:", segmentId)
			segment.file.Close()
			err := os.Remove(segment.path)
			if err != nil {
				logger.Error(fmt.Errorf("DeleteOldSegments: failed to delete old segment file %s: %w", segment.path, err))
				continue
			}
			delete(sm.segmentMap, segmentId)
		} else {
			logger.Warn("Segment with id:", segmentId, "does not exist in the segment map")
		}
	}
}
