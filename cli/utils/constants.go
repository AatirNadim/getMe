package utils

const (
	BaseUrl         = "http://unix"
	GetRoute        = "/get"
	PutRoute        = "/put"
	DeleteRoute     = "/delete"
	ClearStoreRoute = "/clearStore"
	BatchPutRoute   = "/batch-put"
)

type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const (
	SocketPath                 = "/tmp/getMeStore/sockDir/getMe.sock"
	MaxJSONFileSizeBytes int64 = 5 * 1024 * 1024 // 5 MiB
)
