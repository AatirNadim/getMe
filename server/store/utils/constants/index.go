package constants


const (
	SocketPath = "/tmp/getMeStore/getMe.sock"
	StoreDirName   = ".getMeStore"
	DefaultMaxSegmentSize = 50
	MaxEntriesPerSegment = 10000
	ThresholdForCompaction = 10 // if there are more than 10 segments, we trigger compaction
	TotalSegmentsToCompactAtOnce = 5 // we compact 5 segments at a time
)