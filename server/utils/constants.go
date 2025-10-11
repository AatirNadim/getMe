package utils

const (
	GetRoute        = "/get"
	PutRoute        = "/put"
	DeleteRoute     = "/delete"
	ClearStoreRoute = "/clearStore"
	BatchSetRoute   = "/batch-set"
)

type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
