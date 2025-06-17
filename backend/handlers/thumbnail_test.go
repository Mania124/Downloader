package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetThumbnail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock YTDLPCommand function
	YTDLPCommand = func(ctx context.Context, url string) ([]byte, error) {
		// Return a fake thumbnail JSON
		return []byte(`{"thumbnail": "http://example.com/thumb.jpg"}`), nil
	}

	router := gin.Default()
	router.POST("/thumbnail", GetThumbnail)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request",
			requestBody:    map[string]string{"url": "http://example.com/video"},
			expectedStatus: http.StatusOK,
			expectedBody:   `"thumbnail":"http://example.com/thumb.jpg"`,
		},
		{
			name:           "Missing URL",
			requestBody:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Invalid or missing URL"`,
		},
		{
			name:           "Invalid URL",
			requestBody:    map[string]string{"url": "not-a-url"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Invalid or missing URL"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/thumbnail", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
			if !bytes.Contains(resp.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("expected body to contain %s, got %s", tt.expectedBody, resp.Body.String())
			}
		})
	}
}
