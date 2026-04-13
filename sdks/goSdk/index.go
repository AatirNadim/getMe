package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AatirNadim/getMe/sdks/goSdk/core"
	"github.com/AatirNadim/getMe/commons"
	"github.com/joho/godotenv"
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", commons.BaseUrl, commons.GetRoute), nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get key: %s, status code: %d", key, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
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

	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload commons.BatchGetRequestBody

	// unmarshalling the json file followed by marshalling it again might seem redundant, but it allows us to validate the structure of the JSON file and ensure that it contains the expected "keys" field. It also allows us to easily convert the JSON data into the format required for the batch get request.

	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	jsonPayload, err := json.Marshal(commons.BatchGetRequestBody{
		Keys: payload.Keys,
	})
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchGetRoute), readerPayload)
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchGetResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchGetResult{}, fmt.Errorf("failed to batch get keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	var batchGetResponse commons.BatchGetResult

	if err := json.Unmarshal(bodyBytes, &batchGetResponse); err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to unmarshal batch get response: %w", err)
	}

	return batchGetResponse, nil
}

func (client *GetMeClient) BatchGetForPayload(payload commons.BatchGetRequestBody) (commons.BatchGetResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchGetRoute), readerPayload)
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchGetResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchGetResult{}, fmt.Errorf("failed to batch get keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchGetResult{}, err
	}

	var batchGetResponse commons.BatchGetResult

	if err := json.Unmarshal(bodyBytes, &batchGetResponse); err != nil {
		return commons.BatchGetResult{}, fmt.Errorf("failed to unmarshal batch get response: %w", err)
	}

	return batchGetResponse, nil
}

func (client *GetMeClient) Put(key, value string) error {

	fmt.Println("Preparing JSON payload for PUT request with key:", key, " and value:", value)
	jsonPayload, err := json.Marshal(commons.PutRequestBody{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	fmt.Println("preparing io reader payload with:", jsonPayload)
	readerPayload := bytes.NewReader(jsonPayload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.PutRoute), readerPayload)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", key)
	q.Add("value", value)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put key: %s, status code: %d", key, resp.StatusCode)
	}

	return nil
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
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload []commons.KeyValue

	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchPutRoute), readerPayload)
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchPutResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchPutResult{}, fmt.Errorf("failed to batch put keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	var batchPutResponse commons.BatchPutResult

	if err := json.Unmarshal(bodyBytes, &batchPutResponse); err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to unmarshal batch put response: %w", err)
	}

	return batchPutResponse, nil
}

func (client *GetMeClient) BatchPutForPayload(payload []commons.KeyValue) (commons.BatchPutResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchPutRoute), readerPayload)
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchPutResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchPutResult{}, fmt.Errorf("failed to batch put keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchPutResult{}, err
	}

	var batchPutResponse commons.BatchPutResult

	if err := json.Unmarshal(bodyBytes, &batchPutResponse); err != nil {
		return commons.BatchPutResult{}, fmt.Errorf("failed to unmarshal batch put response: %w", err)
	}

	return batchPutResponse, nil
}

func (client *GetMeClient) Delete(key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", commons.BaseUrl, commons.DeleteRoute), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete key: %s, status code: %d", key, resp.StatusCode)
	}

	return nil
}

func (client *GetMeClient) BatchDelete(jsonPath string) (commons.BatchDeleteResult, error) {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var payload commons.BatchDeleteRequestBody

	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to unmarshal JSON file: %w", err)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchDeleteRoute), readerPayload)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to batch delete keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	var batchDeleteResponse commons.BatchDeleteResult

	if err := json.Unmarshal(bodyBytes, &batchDeleteResponse); err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to unmarshal batch delete response: %w", err)
	}

	return batchDeleteResponse, nil
}

func (client *GetMeClient) BatchDeleteForPayload(payload commons.BatchDeleteRequestBody) (commons.BatchDeleteResult, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.BatchDeleteRoute), readerPayload)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to batch delete keys, status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return commons.BatchDeleteResult{}, err
	}

	var batchDeleteResponse commons.BatchDeleteResult

	if err := json.Unmarshal(bodyBytes, &batchDeleteResponse); err != nil {
		return commons.BatchDeleteResult{}, fmt.Errorf("failed to unmarshal batch delete response: %w", err)
	}

	return batchDeleteResponse, nil
}

func (client *GetMeClient) ClearStore() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", commons.BaseUrl, commons.ClearStoreRoute), nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to clear store, status code: %d", resp.StatusCode)
	}

	return nil
}
