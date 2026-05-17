package commons

const (
	SocketPath                 = "/tmp/getMeStore/sockDir/getMe.sock"
	BaseUrl                    = "http://unix"
	GetRoute                   = "/get"
	PutRoute                   = "/put"
	DeleteRoute                = "/delete"
	ClearStoreRoute            = "/clearStore"
	BatchPutRoute              = "/batch-put"
	BatchGetRoute              = "/batch-get"
	BatchDeleteRoute           = "/batch-delete"
	MaxJSONFileSizeBytes int64 = 5 * 1024 * 1024 // 5 MiB
)
