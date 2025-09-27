package store

import (
	"fmt"
	"getMeMod/store/logger"
	"os"
	"path/filepath"
	"sync"
)



type SegmentManager struct {
	mu sync.RWMutex
	basePath string
	segmentMap map[uint32]*Segment
	// stores the index of the next segment to be created
	activeId uint32
}


func NewSegmentManager(basePath string, centralHashTable *HashTable) (*SegmentManager, error) {


	sm := &SegmentManager{
		segmentMap: make(map[uint32]*Segment),
		basePath:   basePath,
		activeId: 0,
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
		if _, err := sm.CreateNewSegment(basePath); err != nil {
			return nil, fmt.Errorf("failed to create active segment: %w", err)
		}
	}

	return sm, nil
}



func (sm *SegmentManager) populateSegmentMap(basePath string, centralHashTable *HashTable) error {

	logger.Info("Segments already exist, loading them from the disk to the current kv instance...")

	// find all segment files in the base path
	paths, err := filepath.Glob(filepath.Join(basePath, "segment_*.log"))
	if err != nil {
		logger.Error("failed to list segment files basePath", basePath, "error", err)
		return fmt.Errorf("failed to list segment files in %s: %w", basePath, err)
	}
	if paths == nil {
		logger.Error("no segments found in " + basePath)
		return nil // No segments found is not an error
	}

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

		sm.activeId = max(sm.activeId, id)

		// assign the segment mapped to its id
		sm.segmentMap[uint32(id)] = segment

		wg.Add(1)
		go func(seg *Segment) {
			defer wg.Done()
			ht, err := seg.ReadAllEntries()
			if err != nil {
				logger.Error("Failed to read entries from segment %d: %v", seg.id, err)
				return
			}
			ch <- ht
		}(segment)
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


	// if the latest entry is a deletion entry, simply remove it from the hash table
	centralHashTable.DeleteDeletionEntries()

	logger.Info("loaded segments from the disk")
	// increment the activeId to be one more than the max id found
	sm.activeId += 1
	return nil

}


// func (sm *SegmentManager) ReadAllSegments() (*HashTable, error) {
// 	var wg sync.WaitGroup
// 	ch := make(chan *HashTable, len(sm.segmentMap))

// 	for _, segment := range sm.segmentMap {
// 		wg.Add(1)
// 		go func(seg *Segment) {
// 			defer wg.Done()
// 			ht, err := seg.ReadAllEntries()
// 			if err != nil {
// 				logger.Error("Failed to read entries from segment %d: %v", seg.id, err)
// 				return
// 			}
// 			ch <- ht
// 		}(segment)
// 	}

// 	// a goroutine to close the channel when all segment reads are done
// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 	}()

// 	// Reducer
// 	finalHt := NewHashTable()
// 	for ht := range ch {
// 		finalHt.Merge(ht)
// 	}

// 	return finalHt, nil
// }

// create a new segment, append it to the segment list and return it
func (sm *SegmentManager) CreateNewSegment(basePath string) (* Segment, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()


	logger.Info("Creating a new segment with id:", sm.activeId)

	// create a new segment with id = sm.activeId
	segment, err := NewSegment(sm.activeId, basePath)
	if err != nil {
		return nil, err
	}

	// add the new segment to the segment map and increment the activeId
	sm.segmentMap[sm.activeId] = segment
	// increment the activeId for the next segment
	sm.activeId += 1


	return segment, nil
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
	return sm.activeId, offset, nil
}


// reads an entry from a specific segment at a specific offset and returns it along with the offset for the next entry
func (sm *SegmentManager) Read(segmentId uint32, offset uint32) (*Entry, uint32, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	logger.Info("segment manager: Reading entry from segment", segmentId, "at offset", offset)

	if segmentId >= sm.activeId - 1 {
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

func (sm *SegmentManager) Delete(key []byte) (uint32, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Info("segment manager: Deleting entry with key:", string(key))

	deletionEntry, deletionEntryCreationErr := CreateDeletionEntry(key)

	if deletionEntryCreationErr != nil {
		logger.Error("segment manager: failed to create deletion entry:", deletionEntryCreationErr)
		return 0, deletionEntryCreationErr
	}

	offset, err := sm.appendEntryToLatestSegment(deletionEntry)

	// deletionEntry, deletionEntryCreationErr := activeSegment.CreateDeletionEntry(key)

	// if deletionEntryCreationErr != nil {
	// 	return 0, deletionEntryCreationErr
	// }

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

	sm.activeId = 0
}



func (sm *SegmentManager) appendEntryToLatestSegment(entry *Entry) (uint32, error) {


	// active id will always hold the id of the next segment to be created
	currentSegment := sm.segmentMap[sm.activeId - 1]

	if !isSpaceAvailableInCurrentSegment(currentSegment, entry) {
		newSegment, newSegmentCreationError := sm.CreateNewSegment(sm.basePath)

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