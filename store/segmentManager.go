package store

import (
	"fmt"
	"getMeMod/store/utils"
	"path/filepath"
	"sync"
)



type SegmentManager struct {
	mu sync.RWMutex
	basePath string
	segments []*Segment
	activeId int
}


func NewSegmentManager() *SegmentManager {
	return &SegmentManager{
		segments: make([]*Segment, 0),
	}
}



// loads existing segments from the disk, from the base path
func loadSegments(basePath string) (*SegmentManager, error) {
	paths, err := filepath.Glob(filepath.Join(basePath, "segment_*.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to list segment files in %s: %w", basePath, err)
	}
	if paths == nil {
		return nil, fmt.Errorf("no segments found in %s", basePath)
	}
	segments := make([]*Segment, 0, len(paths))

	segManager := &SegmentManager{}

	for _, path := range paths {
		var id int
		_, err := fmt.Sscanf(filepath.Base(path), "segment_%d.log", &id)
		if err != nil {
			return nil, err
		}

		segment, err := OpenSegment(id, basePath)
		if err != nil {
			return nil, err
		}

		segments[id] = segment // plugging the segments based on the index present in their names

	
		if id > segManager.activeId {
			segManager.activeId = id + 1
		}
	}

	return &SegmentManager{
		segments: segments, // this list is sorted based on the index
		activeId: segManager.activeId,
	}, nil

}


// create a new segment, append it to the segment list and return it
func (sm *SegmentManager) CreateNewSegment(basePath string) (*Segment, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// create a new segment with id = len(segments)
	segment, err := NewSegment(len(sm.segments), basePath)
	if err != nil {
		return nil, err
	}

	sm.segments = append(sm.segments, segment)
	

	return sm.segments[len(sm.segments) - 1], nil
}

func (sm *SegmentManager) Append(entry *Entry) (uint32, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	activeSegment := sm.segments[len(sm.segments) - 1]


	_, err := activeSegment.Append(entry)


	if err != nil {
		if err == utils.ErrSegmentFull {
			newSegment, err := sm.CreateNewSegment(sm.basePath)

			if err != nil {
				return 0, err
			}

			res, err := newSegment.Append(entry)
			if err != nil {
				return 0, err
			}

			return res, nil
		} else {
			return 0, err
		}
	}
	return 0, nil
}


// reads an entry from a specific segment at a specific offset
func (sm *SegmentManager) Read(segmentId int, offset uint32) (*Entry, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if segmentId >= len(sm.segments) {
		return nil, fmt.Errorf("segment %d does not exist", segmentId)
	}

	segment := sm.segments[segmentId]
	return segment.Get(offset)
}

func (sm *SegmentManager) Update(entry *Entry) (uint32, int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	activeSegment := sm.segments[len(sm.segments) - 1]

	res, err := activeSegment.Append(entry)
	if err != nil {
		if err == utils.ErrSegmentFull {
			newSegment, err := sm.CreateNewSegment(sm.basePath)

			if err != nil {
				return 0, 0, err
			}

			res, err := newSegment.Append(entry)
			if err != nil {
				return 0, 0, err
			}

			return res, newSegment.id, nil
		} else {
			return 0, 0, err
		}
	}
	return res, activeSegment.id, nil
}

func (sm *SegmentManager) Delete(key []byte) (uint32, int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	activeSegment := sm.segments[len(sm.segments) - 1]

	deletionEntry := CreateDeletionEntry(key)

	_, err := sm.appendEntryToLatestSegment(deletionEntry)

	if err != nil {
		return 0, 0, err
	}

	return 0, activeSegment.id, nil	

}




func (sm *SegmentManager) appendEntryToLatestSegment(entry *Entry) (uint32, error) {
	activeSegment := sm.segments[len(sm.segments) - 1]

	res, err := activeSegment.Append(entry)
	if err != nil {
		if err == utils.ErrSegmentFull {
			newSegment, err := sm.CreateNewSegment(sm.basePath)

			if err != nil {
				return 0, err
			}

			res, err := newSegment.Append(entry)
			if err != nil {
				return 0, err
			}

			return res, nil
		} else {
			return 0, err
		}
	}
	return res, nil
}