package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupDownloadRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/download", DownloadVideo)
	return router
}

func TestDownloadVideo_InvalidRequest(t *testing.T) {
	router := setupDownloadRouter()

	// Test with empty JSON body (invalid)
	reqBody := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest(http.MethodPost, "/download", reqBody)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestDownloadVideo_InvalidFormat(t *testing.T) {
	router := setupDownloadRouter()

	// Test with valid URL but invalid format
	reqBody := bytes.NewBufferString(`{"url":"https://example.com","format":"invalid"}`)
	req, _ := http.NewRequest(http.MethodPost, "/download", reqBody)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
