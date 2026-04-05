package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AatirNadim/getMe/cli/utils"
	"github.com/AatirNadim/getMe/utils/logger"
	"net/http"
	"os"
)

type ServiceLayer struct {
	HttpClient *http.Client
}

func (s *ServiceLayer) GetService(key string) (string, error) {

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodGet,
		URL:    utils.BaseUrl,
		Path:   utils.GetRoute,
		QueryParams: map[string]string{
			"key": key,
		},
	})

	if err != nil {
		logger.Error("Error occurred while creating GET request:", err)
		return "", fmt.Errorf("failed to create GET request for key '%s': %w", key, err)
	}

	respStr, err := utils.ExecuteHTTPRequest(s.HttpClient, req)

	return respStr, nil
}

func (s *ServiceLayer) GetJsonValueService(key string) ([]byte, error) {
	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodGet,
		URL:    utils.BaseUrl,
		Path:   utils.GetRoute,
		QueryParams: map[string]string{
			"key": key,
		},
	})

	if err != nil {
		logger.Error("Error occurred while creating GET request:", err)
		return nil, fmt.Errorf("failed to create GET request for key '%s': %w", key, err)
	}

	resp, err := utils.ExecuteHTTPRequestAndReturnBuffer(s.HttpClient, req)

	return resp, nil
}

func (s *ServiceLayer) BatchGetService(jsonFilePath string) (utils.BatchGetResult, error) {

	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to read file '%s': %w", jsonFilePath, err)
	}

	// being able to parse the JSON input file in the desired format is important!
	var batchGetReq utils.BatchGetRequestBody
	if err := json.Unmarshal(fileContent, &batchGetReq); err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to parse JSON file '%s': %w", jsonFilePath, err)
	}

	jsonPayload, err := json.Marshal(utils.BatchGetRequestBody{
		Keys: batchGetReq.Keys,
	})
	if err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to marshal batch get keys into JSON: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodPost,
		URL:    utils.BaseUrl,
		Path:   utils.BatchGetRoute,
		Body:   readerPayload,
	})
	if err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to create batch get request: %w", err)
	}

	respStr, err := utils.ExecuteHTTPRequest(s.HttpClient, req)
	if err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to execute batch get request: %w", err)
	}

	var result utils.BatchGetResult
	if err := json.Unmarshal([]byte(respStr), &result); err != nil {
		return utils.BatchGetResult{}, fmt.Errorf("failed to unmarshal batch get response: %w", err)
	}

	return result, nil
}

func (s *ServiceLayer) PutService(key, value string) error {

	logger.Debug("Preparing JSON payload for PUT request with key:", key, " and value:", value)
	jsonPayload, err := json.Marshal(utils.PutRequestBody{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	logger.Debug("preparing io reader payload with:", jsonPayload)
	readerPayload := bytes.NewReader(jsonPayload)

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodPut,
		URL:    utils.BaseUrl,
		Path:   utils.PutRoute,
		Body:   readerPayload,
	})

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, err = utils.ExecuteHTTPRequest(s.HttpClient, req)

	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	return nil

}

func (s *ServiceLayer) BatchPutService(keyValuePairs map[string]string) error {

	jsonPayload, err := json.Marshal(keyValuePairs)
	if err != nil {
		return fmt.Errorf("failed to marshal batch data into JSON: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodPost,
		URL:    utils.BaseUrl,
		Path:   utils.BatchPutRoute,
		Body:   readerPayload,
	})

	if err != nil {
		return fmt.Errorf("failed to create batch PUT request: %w", err)
	}

	_, err = utils.ExecuteHTTPRequest(s.HttpClient, req)

	if err != nil {
		return fmt.Errorf("failed to perform batch PUT request: %w", err)
	}

	return nil
}

func (s *ServiceLayer) DeleteService(key string) error {

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodDelete,
		URL:    utils.BaseUrl,
		Path:   utils.DeleteRoute,
		QueryParams: map[string]string{
			"key": key,
		},
	})

	if err != nil {
		logger.Error("Error occurred while creating DELETE request:", err)
		return fmt.Errorf("failed to create DELETE request for key '%s': %w", key, err)
	}

	_, err = utils.ExecuteHTTPRequest(s.HttpClient, req)

	if err != nil {
		logger.Error("Error occurred while executing DELETE request:", err)
		return fmt.Errorf("failed to perform DELETE request for key '%s': %w", key, err)
	}

	return nil
}

func (s *ServiceLayer) BatchDeleteService(batchDeleteReq utils.BatchGetRequestBody) error {
	jsonPayload, err := json.Marshal(batchDeleteReq)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	readerPayload := bytes.NewReader(jsonPayload)

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodPost,
		URL:    utils.BaseUrl,
		Path:   utils.BatchDeleteRoute,
		Body:   readerPayload,
	})
	if err != nil {
		return fmt.Errorf("failed to create batch delete request: %w", err)
	}

	_, err = utils.ExecuteHTTPRequest(s.HttpClient, req)
	if err != nil {
		return fmt.Errorf("failed to execute batch delete request: %w", err)
	}

	return nil

}

func (s *ServiceLayer) ClearStoreService() error {

	req, err := utils.CreateHTTPRequest(utils.RequestOptions{
		Method: http.MethodPost,
		URL:    utils.BaseUrl,
		Path:   utils.ClearStoreRoute,
	})

	if err != nil {
		logger.Error("Error occurred while creating Clear Store request:", err)
		return fmt.Errorf("failed to create Clear Store request: %w", err)
	}

	_, err = utils.ExecuteHTTPRequest(s.HttpClient, req)

	if err != nil {
		logger.Error("Error occurred while executing Clear Store request:", err)
		return fmt.Errorf("failed to perform Clear Store request: %w", err)
	}

	return nil
}
