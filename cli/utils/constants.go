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
	SocketPath = "/tmp/getMeStore/getMe.sock"
)
