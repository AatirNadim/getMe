package constants

const (
	SocketPath                   = "/tmp/getMeStore/sockDir/getMe.sock"
	LogsDirName                  = "/tmp/getMeStore/dumpDir/index.log"
	StoreDirName                 = "/var/lib/getMeStore/dataDir/segments"
	CompactedStoreDirName        = "/var/lib/getMeStore/dataDir/compactedSegments"
	DefaultMaxSegmentSize        = 100 * 1024 // 100KB
	MaxEntriesPerSegment         = 10000
	ThresholdForCompaction       = 10        // if there are more than 10 segments, we trigger compaction
	TotalSegmentsToCompactAtOnce = 5         // we compact 5 segments at a time
	MaxChunkSize                 = 64 * 1024 // 64KB
)
