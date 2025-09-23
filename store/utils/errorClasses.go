package utils

import "errors"


var (	ErrKeyNotFound = errors.New("key not found")
	ErrKeyDeleted  = errors.New("key has been deleted")
	ErrInvalidEntry = errors.New("invalid entry")
	ErrSegmentFull = errors.New("segment is full")
)