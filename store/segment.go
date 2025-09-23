package store

import "sync"


const DefaultMaxSegmentSize = 1024 * 1024 * 10 // 10 MB

const MaxEntriesPerSegment = 10000


type Segment struct {
	mu sync.RWMutex
	id int
	path string
	entryCount int
	size int
	isActive bool
}


func (segment *Segment) Append(entry* Entry) (uint32, error) {
	segment.mu.Lock()

	defer segment.mu.Unlock()

	if segment.entryCount >= MaxEntriesPerSegment || segment.size + entry.Size() > DefaultMaxSegmentSize {
		return 0, ErrSegmentFull
	}

	data := entry.Serialize()

	

}