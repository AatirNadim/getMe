package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AatirNadim/getMe/commons"

	
	"github.com/joho/godotenv"
	"github.com/AatirNadim/getMe/sdks/goSdk/core"
)

type GetMeClient struct {
	httpClient *http.Client
}

func (client *GetMeClient) Init() error {

	err := godotenv.Load()
	if err != nil {
		return err
	}
	var socketPath string
	socketPath = os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = commons.SocketPath
	}

	httpClient, err := core.CreateHttpClient(socketPath)
	if err != nil {
		return err
	}
	client.httpClient = httpClient
	return nil
}

func (client *GetMeClient) Get(key string) (string, error) {
	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodGet,
		URL:    commons.BaseUrl,
		Path:   commons.GetRoute,
		QueryParams: map[string]string{
			"key": key,
		},
	})
	if err != nil {
		return "", err
	}

	respStr, err := commons.ExecuteHTTPRequest(client.httpClient, req)
	if err != nil {
		return "", err
	}

	return respStr, nil
}

// expects the stored value to be a valid JSON document.
func (client *GetMeClient) GetJSON(key string, out interface{}) error {
	value, err := client.Get(key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(value), out); err != nil {
		return fmt.Errorf("value for key '%s' is not valid JSON: %w", key, err)
	}
	return nil
}

func (client *GetMeClient) BatchGet(jsonPath string) (commons.BatchGetResult, error) {

	fileContent, err := os.ReadFile(jsonPath)
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload commons.BatchGetRequestBody
	if err := json.Unmarshal(fileContent, &payload); err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	return client.BatchGetForPayload(commons.BatchGetRequestBody{
		Keys: payload.Keys,
	})
}

func (client *GetMeClient) BatchGetForPayload(payload commons.BatchGetRequestBody) (commons.BatchGetResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodPost,
		URL:    commons.BaseUrl,
		Path:   commons.BatchGetRoute,
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	respStr, err := commons.ExecuteHTTPRequest(client.httpClient, req)
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	var batchGetResponse commons.BatchGetResult
	if err := json.Unmarshal([]byte(respStr), &batchGetResponse); err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to unmarshal batch get response: %w", err)
	}

	return batchGetResponse, nil
}

func (client *GetMeClient) Put(key, value string) error {
	jsonPayload, err := json.Marshal(commons.PutRequestBody{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodPost,
		URL:    commons.BaseUrl,
		Path:   commons.PutRoute,
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return err
	}

	_, err = commons.ExecuteHTTPRequest(client.httpClient, req)
	return err
}

// PutJSON marshals v as JSON and stores it as the value for key.
func (client *GetMeClient) PutJSON(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON value: %w", err)
	}
	return client.Put(key, string(data))
}

func (client *GetMeClient) BatchPut(jsonPath string) (commons.BatchPutResult, error) {
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload []commons.KeyValue
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	return client.BatchPutForPayload(payload)
}

func (client *GetMeClient) BatchPutForPayload(payload []commons.KeyValue) (commons.BatchPutResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodPost,
		URL:    commons.BaseUrl,
		Path:   commons.BatchPutRoute,
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	respStr, err := commons.ExecuteHTTPRequest(client.httpClient, req)
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	var batchPutResponse commons.BatchPutResult
	if err := json.Unmarshal([]byte(respStr), &batchPutResponse); err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to unmarshal batch put response: %w", err)
	}

	return batchPutResponse, nil
}

func (client *GetMeClient) Delete(key string) error {
	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodDelete,
		URL:    commons.BaseUrl,
		Path:   commons.DeleteRoute,
		QueryParams: map[string]string{
			"key": key,
		},
	})
	if err != nil {
		return err
	}

	_, err = commons.ExecuteHTTPRequest(client.httpClient, req)
	return err
}

func (client *GetMeClient) BatchDelete(jsonPath string) (commons.BatchDeleteResult, error) {
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload commons.BatchDeleteRequestBody
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	return client.BatchDeleteForPayload(payload)
}

func (client *GetMeClient) BatchDeleteForPayload(payload commons.BatchDeleteRequestBody) (commons.BatchDeleteResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodDelete,
		URL:    commons.BaseUrl,
		Path:   commons.BatchDeleteRoute,
		Body:   bytes.NewReader(jsonPayload),
	})
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	respStr, err := commons.ExecuteHTTPRequest(client.httpClient, req)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	var batchDeleteResponse commons.BatchDeleteResult
	if err := json.Unmarshal([]byte(respStr), &batchDeleteResponse); err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to unmarshal batch delete response: %w", err)
	}

	return batchDeleteResponse, nil
}

func (client *GetMeClient) ClearStore() error {
	req, err := commons.CreateHTTPRequest(commons.RequestOptions{
		Method: http.MethodDelete,
		URL:    commons.BaseUrl,
		Path:   commons.ClearStoreRoute,
	})
	if err != nil {
		return err
	}

	_, err = commons.ExecuteHTTPRequest(client.httpClient, req)
	return err
}
