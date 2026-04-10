package correctnessCheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AatirNadim/getMe/tests/util"
)

func TestPutController_EdgeCases(t *testing.T) {
	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	handler := controllers.PutController()

	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Method Not Allowed (GET instead of POST)",
			method:         http.MethodGet,
			url:            "/put",
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed",
		},
		{
			name:           "Invalid JSON format",
			method:         http.MethodPost,
			url:            "/put",
			requestBody:    `{"key": "test_key",}`, // syntax error
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to parse JSON body",
		},
		{
			name:           "Missing Key Parameter",
			method:         http.MethodPost,
			url:            "/put",
			requestBody:    map[string]string{"value": "test_value"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing key or value parameter",
		},
		{
			name:           "Missing Value Parameter",
			method:         http.MethodPost,
			url:            "/put",
			requestBody:    map[string]string{"key": "test_key"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing key or value parameter",
		},
		{
			name:           "Missing Both Key and Value",
			method:         http.MethodPost,
			url:            "/put",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing key or value parameter",
		},
		{
			name:           "Valid Payload",
			method:         http.MethodPost,
			url:            "/put",
			requestBody:    map[string]string{"key": "test_key", "value": "test_value"},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf("Successfully set value for key '%s'", "test_key"),
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
				// Ensure nil becomes an empty body instead of literal "null"
				bodyBytes = nil
			default:
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(tc.method, tc.url, bytes.NewReader(bodyBytes))
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
