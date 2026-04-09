package correctnessCheck

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AatirNadim/getMe/benchmarking/util"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
)

func TestBatchGetController_EdgeCases(t *testing.T) {

	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	err := controllers.StoreInstance.Put("valid_key_1", "value1")
	if err != nil {
		t.Fatalf("Failed to set up store with valid key: %v", err)
	}

	handler := controllers.BatchGetController()

	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Method Not Allowed (GET instead of POST)",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed",
		},
		{
			name:           "Invalid JSON format",
			method:         http.MethodPost,
			requestBody:    `{"keys": ["key1",}`, // syntax error
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to parse JSON body",
		},
		{
			name:           "Empty Keys List",
			method:         http.MethodPost,
			requestBody:    map[string][]string{"keys": {}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "empty keys list",
		},
		{
			name:   "Empty Keys List after Duplicate Removal",
			method: http.MethodPost,
			// Testing if delete duplicate leaves an empty list (edge case theoretically if it was ["", ...] but payload.Keys length 0 handles it)
			requestBody:    map[string][]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "empty keys list",
		},
		{
			name:           "Too Many Keys (> 10000)",
			method:         http.MethodPost,
			requestBody:    map[string][]string{"keys": util.GenerateDummyKeys(10001)},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "too many keys",
		},
		{
			name:   "Valid Keys with Duplicates and Missing keys",
			method: http.MethodPost,
			// Payload includes a duplicate key, an existing key, and a missing key
			requestBody:    map[string][]string{"keys": {"valid_key_1", "valid_key_1", "missing_key"}},
			expectedStatus: http.StatusOK,
			// Should return `{"found":{"valid_key_1":"value1"},"errors":{},"notFound":["missing_key"]}` or similar based on structure
			expectedBody: `"valid_key_1":"value1"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error

			switch v := tc.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			case nil:
				bodyBytes = nil
			default:
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(tc.method, "/batch-get", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.expectedStatus, rr.Code, rr.Body.String())
			}

			if tc.expectedBody != "" && !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("Expected body to contain %q, but got %q", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestBatchGetController_MaxBodySize(t *testing.T) {
	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	handler := controllers.BatchGetController()

	// Create a payload larger than MaxBodySize constraint
	// constants.MaxBodySize is typically around 1MB or similar
	largeSize := constants.MaxBodySize + 100
	largeBytes := make([]byte, largeSize)
	for i := range largeBytes {
		largeBytes[i] = 'A' // Fill with arbitrary large text
	}

	req := httptest.NewRequest(http.MethodPost, "/batch-get", bytes.NewReader(largeBytes))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// MaxBytesReader should throw an error when attempting to read the oversized body
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected BAD REQUEST status %d for entity too large, got %d. Response: %v", http.StatusBadRequest, rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "failed to read request body") {
		t.Errorf("Expected 'failed to read request body', got %s", rr.Body.String())
	}
}
