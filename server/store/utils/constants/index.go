package constants

const (
	SocketPath                   = "/tmp/getMeStore/getMe.sock"
	StoreDirName                 = ".getMeStore/data"
	CompactedStoreDirName        = ".getMeStore/compacted_data"
	DefaultMaxSegmentSize        = 100 // in bytes, 100 bytes for testing, should be in MBs in real world
	MaxEntriesPerSegment         = 10000
	ThresholdForCompaction       = 10        // if there are more than 10 segments, we trigger compaction
	TotalSegmentsToCompactAtOnce = 5         // we compact 5 segments at a time
	MaxChunkSize                 = 64 * 1024 // 64KB
)
