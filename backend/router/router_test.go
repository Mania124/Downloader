package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouterForTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	SetupRoutes(r)
	return r
}

func TestHealthCheckRoute(t *testing.T) {
	r := setupRouterForTest()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
}

func TestCORSConfig(t *testing.T) {
	os.Setenv("FRONTEND_ORIGIN", "http://test-origin.com")
	r := gin.Default()
	SetupRoutes(r)

	req, _ := http.NewRequest(http.MethodOptions, "/download", nil)
	req.Header.Set("Origin", "http://test-origin.com")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	allowedOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowedOrigin != "http://test-origin.com" {
		t.Errorf("Expected CORS header to allow 'http://test-origin.com', got '%s'", allowedOrigin)
	}
}
