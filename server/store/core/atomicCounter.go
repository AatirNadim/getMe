package core


import (
	"sync/atomic"
)


// this atomic counter is going to be used to get the segment id for the next segment
// to be used for aquiring new segment ids as well, (for compaction), so it needs to be thread safe
type AtomicCounter struct {
	value uint32
}

func NewAtomicCounter(initialValue uint32) *AtomicCounter {
	return &AtomicCounter{
		value: initialValue,
	}
}


func (ac *AtomicCounter) Set(value uint32) {
	atomic.StoreUint32(&ac.value, value)
}

// Increments the counter by 1 and returns the value
func (ac *AtomicCounter) Next() uint32 {
  return atomic.AddUint32(&ac.value, 1)
}


// Reserve reserves n values and returns the first value in the reservation
// the atomic counter now holds the first value beyond the reserved range
// This operation is going to be used to reserve segment ids for compaction
func (ac *AtomicCounter) Reserve(n uint32) uint32 {
  return atomic.AddUint32(&ac.value, n) - n
}

// Get returns the current value of the counter without incrementing it
func (ac *AtomicCounter) Get() uint32 {
	return atomic.LoadUint32(&ac.value)
}

	