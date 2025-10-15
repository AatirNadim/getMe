package utils

const (
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