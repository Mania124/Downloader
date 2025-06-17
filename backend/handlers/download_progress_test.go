package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupDownloadProgressRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/download/stream", DownloadWithProgress)
	return router
}

func TestDownloadWithProgress_MissingURL(t *testing.T) {
	router := setupDownloadProgressRouter()

	req, _ := http.NewRequest(http.MethodGet, "/download/stream?format=video", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestDownloadWithProgress_InvalidFormat(t *testing.T) {
	router := setupDownloadProgressRouter()

	req, _ := http.NewRequest(http.MethodGet, "/download/stream?url=https://example.com&format=invalid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}
