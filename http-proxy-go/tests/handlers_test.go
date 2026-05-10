package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AatirNadim/getMe/commons"
	"github.com/AatirNadim/getMe/http-proxy-go/handlers"
)

// MockClient implements handlers.GetMeClientInterface
type MockClient struct {
	MockGet                   func(key string) (string, error)
	MockPut                   func(key, value string) error
	MockDelete                func(key string) error
	MockBatchGetForPayload    func(payload commons.BatchGetRequestBody) (commons.BatchGetResult, error)
	MockBatchPutForPayload    func(payload map[string]string) (commons.BatchPutResult, error)
	MockBatchDeleteForPayload func(payload commons.BatchDeleteRequestBody) (commons.BatchDeleteResult, error)
	MockClearStore            func() error
}

func (m *MockClient) Get(key string) (string, error) {
	return m.MockGet(key)
}

func (m *MockClient) Put(key, value string) error {
	return m.MockPut(key, value)
}

func (m *MockClient) Delete(key string) error {
	return m.MockDelete(key)
}

func (m *MockClient) BatchGetForPayload(payload commons.BatchGetRequestBody) (commons.BatchGetResult, error) {
	return m.MockBatchGetForPayload(payload)
}

func (m *MockClient) BatchPutForPayload(payload map[string]string) (commons.BatchPutResult, error) {
	return m.MockBatchPutForPayload(payload)
}

func (m *MockClient) BatchDeleteForPayload(payload commons.BatchDeleteRequestBody) (commons.BatchDeleteResult, error) {
	return m.MockBatchDeleteForPayload(payload)
}

func (m *MockClient) ClearStore() error {
	return m.MockClearStore()
}

func TestGetHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockGet = func(key string) (string, error) {
			if key != "test-key" {
				return "", errors.New("unexpected key")
			}
			return "test-value", nil
		}

		req := httptest.NewRequest(http.MethodGet, "/get?key=test-key", nil)
		rr := httptest.NewRecorder()

		proxy.GetHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
		if rr.Body.String() != "test-value" {
			t.Errorf("Expected body test-value, got %s", rr.Body.String())
		}
	})

	t.Run("Missing Key Parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get", nil)
		rr := httptest.NewRecorder()
		proxy.GetHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mockClient.MockGet = func(key string) (string, error) {
			return "", errors.New("network failure")
		}

		req := httptest.NewRequest(http.MethodGet, "/get?key=test-key", nil)
		rr := httptest.NewRecorder()
		proxy.GetHandler(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status InternalServerError, got %d", rr.Code)
		}
	})

	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/get?key=test-key", nil)
		rr := httptest.NewRecorder()
		proxy.GetHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status MethodNotAllowed, got %d", rr.Code)
		}
	})
}

func TestPutHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockPut = func(key, value string) error {
			if key != "k1" || value != "v1" {
				return errors.New("unexpected data")
			}
			return nil
		}

		body, _ := json.Marshal(map[string]string{"key": "k1", "value": "v1"})
		req := httptest.NewRequest(http.MethodPost, "/put", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		proxy.PutHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Missing Fields", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"key": "k1"})
		req := httptest.NewRequest(http.MethodPost, "/put", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		proxy.PutHandler(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})
}

func TestDeleteHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockDelete = func(key string) error {
			if key != "test-key" {
				return errors.New("unexpected key")
			}
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/delete?key=test-key", nil)
		rr := httptest.NewRecorder()

		proxy.DeleteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Missing Key Parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
		rr := httptest.NewRecorder()
		proxy.DeleteHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mockClient.MockDelete = func(key string) error {
			return errors.New("network failure")
		}

		req := httptest.NewRequest(http.MethodDelete, "/delete?key=test-key", nil)
		rr := httptest.NewRecorder()
		proxy.DeleteHandler(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status InternalServerError, got %d", rr.Code)
		}
	})

	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/delete?key=test-key", nil)
		rr := httptest.NewRecorder()
		proxy.DeleteHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status MethodNotAllowed, got %d", rr.Code)
		}
	})
}

func TestBatchGetHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockBatchGetForPayload = func(payload commons.BatchGetRequestBody) (commons.BatchGetResult, error) {
			if len(payload.Keys) != 1 || payload.Keys[0] != "k1" {
				return commons.BatchGetResult{}, errors.New("unexpected payload")
			}
			return commons.BatchGetResult{Found: map[string]string{"k1": "v1"}}, nil
		}

		body, _ := json.Marshal(commons.BatchGetRequestBody{Keys: []string{"k1"}})
		req := httptest.NewRequest(http.MethodPost, "/batchGet", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		proxy.BatchGetHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Empty Keys", func(t *testing.T) {
		body, _ := json.Marshal(commons.BatchGetRequestBody{Keys: []string{}})
		req := httptest.NewRequest(http.MethodPost, "/batchGet", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		proxy.BatchGetHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/batchGet", bytes.NewReader([]byte("invalid json")))
		rr := httptest.NewRecorder()
		proxy.BatchGetHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})
}

func TestBatchPutHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockBatchPutForPayload = func(payload map[string]string) (commons.BatchPutResult, error) {
			if val, ok := payload["k1"]; !ok || val != "v1" {
				return commons.BatchPutResult{}, errors.New("unexpected payload")
			}
			return commons.BatchPutResult{Successful: 1}, nil
		}

		body, _ := json.Marshal(map[string]string{"k1": "v1"})
		req := httptest.NewRequest(http.MethodPost, "/batchPut", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		proxy.BatchPutHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Empty Payload", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{})
		req := httptest.NewRequest(http.MethodPost, "/batchPut", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		proxy.BatchPutHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})
}

func TestBatchDeleteHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request", func(t *testing.T) {
		mockClient.MockBatchDeleteForPayload = func(payload commons.BatchDeleteRequestBody) (commons.BatchDeleteResult, error) {
			if len(payload.Keys) != 1 || payload.Keys[0] != "k1" {
				return commons.BatchDeleteResult{}, errors.New("unexpected payload")
			}
			return commons.BatchDeleteResult{Successful: 1}, nil
		}

		body, _ := json.Marshal(commons.BatchDeleteRequestBody{Keys: []string{"k1"}})
		req := httptest.NewRequest(http.MethodDelete, "/batchDelete", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		proxy.BatchDeleteHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Empty Keys", func(t *testing.T) {
		body, _ := json.Marshal(commons.BatchDeleteRequestBody{Keys: []string{}})
		req := httptest.NewRequest(http.MethodDelete, "/batchDelete", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		proxy.BatchDeleteHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %d", rr.Code)
		}
	})
}

func TestClearStoreHandler(t *testing.T) {
	mockClient := &MockClient{}
	proxy := &handlers.HttpProxy{Client: mockClient}

	t.Run("Valid Request Post", func(t *testing.T) {
		mockClient.MockClearStore = func() error {
			return nil
		}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		rr := httptest.NewRecorder()

		proxy.ClearStoreHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Valid Request Delete", func(t *testing.T) {
		mockClient.MockClearStore = func() error {
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/clear", nil)
		rr := httptest.NewRecorder()

		proxy.ClearStoreHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	t.Run("Client Error", func(t *testing.T) {
		mockClient.MockClearStore = func() error {
			return errors.New("network failure")
		}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		rr := httptest.NewRecorder()
		proxy.ClearStoreHandler(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status InternalServerError, got %d", rr.Code)
		}
	})

	t.Run("Invalid Method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/clear", nil)
		rr := httptest.NewRecorder()
		proxy.ClearStoreHandler(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status MethodNotAllowed, got %d", rr.Code)
		}
	})
}
