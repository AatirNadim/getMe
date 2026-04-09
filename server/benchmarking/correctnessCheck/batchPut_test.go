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

func TestBatchPutController_EdgeCases(t *testing.T) {

	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	handler := controllers.BatchPutController()

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
			requestBody:    `{"key1": "value1",}`, // syntax error
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to parse JSON body",
		},
		{
			name:           "Invalid JSON Type (Array instead of Object)",
			method:         http.MethodPost,
			requestBody:    `[{"key1": "value1"}]`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to parse JSON body",
		},
		{
			name:           "Empty Batch",
			method:         http.MethodPost,
			requestBody:    map[string]string{},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Valid Batch",
			method: http.MethodPost,
			requestBody: map[string]string{
				"test_put_key_1": "val1",
				"test_put_key_2": "val2",
			},
			expectedStatus: http.StatusOK,
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

			req := httptest.NewRequest(tc.method, "/batch-put", bytes.NewReader(bodyBytes))
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

func TestBatchPutController_MaxBodySize(t *testing.T) {
	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	handler := controllers.BatchPutController()

	// Create a payload larger than MaxBodySize constraint
	largeSize := constants.MaxBodySize + 100
	largeBytes := make([]byte, largeSize)
	for i := range largeBytes {
		largeBytes[i] = 'A' // Fill with arbitrary mock text
	}

	req := httptest.NewRequest(http.MethodPost, "/batch-put", bytes.NewReader(largeBytes))
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
