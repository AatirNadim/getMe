package constants

const (
	SocketPath       = "/tmp/getMeStore/sockDir/getMe.sock"
	BaseUrl          = "http://unix"
	GetRoute         = "/get"
	PutRoute         = "/put"
	DeleteRoute      = "/delete"
	ClearStoreRoute  = "/clearStore"
	BatchPutRoute    = "/batch-put"
	BatchGetRoute    = "/batch-get"
	BatchDeleteRoute = "/batch-delete"
)

type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BatchGetRequestBody struct {
	Keys []string `json:"keys"`
}

type BatchGetResult struct {
	Found    map[string]string `json:"found"`    // key is the key that was found, value is the corresponding value
	NotFound []string          `json:"notFound"` // list of keys that were not found in the store
	Errors   map[string]string `json:"errors"`   // key is the key that failed to get, value is the error message
}

type BatchPutResult struct {
	Successful int               `json:"successful"`
	Failed     map[string]string `json:"failed"` // key is the key that failed to put, value is the error message
}
type BatchDeleteRequestBody = BatchGetRequestBody

type BatchDeleteResult = BatchPutResult
