package router

import (
	"downloader/handlers"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Get allowed origin from environment variable, fallback to localhost:5173
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/", handlers.HealthCheck)
	r.POST("/download", handlers.DownloadVideo)
	r.POST("/thumbnail", handlers.GetThumbnail)
	r.GET("/download/stream", handlers.DownloadWithProgress)
}
