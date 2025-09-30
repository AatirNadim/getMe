package utils

import "errors"


var (	
	ErrorInSerialization = errors.New("error in serialization")
	ErrorInDeserialization = errors.New("error in deserialization")
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyDeleted  = errors.New("key has been deleted")
	ErrInvalidEntry = errors.New("invalid entry")
	ErrSegmentFull = errors.New("segment is full")
	ErrFileNotFoundOrNotAccessible = errors.New("file not found or is not accessible")
)