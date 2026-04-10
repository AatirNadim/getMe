package correctnessCheck

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AatirNadim/getMe/tests/util"
)

func TestGetController_EdgeCases(t *testing.T) {
	controllers, cleanup := util.SetupStoreForCorrectnessCheck(t)
	defer cleanup()

	err := controllers.StoreInstance.Put("valid_key_1", "value1")
	if err != nil {
		t.Fatalf("Failed to set up store with valid key: %v", err)
	}

	handler := controllers.GetController()

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Method Not Allowed (POST instead of GET)",
			method:         http.MethodPost,
			url:            "/get?key=valid_key_1",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed",
		},
		{
			name:           "Missing Key Parameter",
			method:         http.MethodGet,
			url:            "/get",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing key parameter",
		},
		{
			name:           "Empty Key Parameter",
			method:         http.MethodGet,
			url:            "/get?key=",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing key parameter",
		},
		{
			name:           "Key Not Found",
			method:         http.MethodGet,
			url:            "/get?key=missing_key",
			expectedStatus: http.StatusNotFound,
			expectedBody:   fmt.Sprintf("key '%s'", "missing_key"),
		},
		{
			name:           "Valid Key Found",
			method:         http.MethodGet,
			url:            "/get?key=valid_key_1",
			expectedStatus: http.StatusOK,
			expectedBody:   "value1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
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
