package utils

const (
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
