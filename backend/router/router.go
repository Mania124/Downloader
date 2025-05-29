package router

import (
	"downloader/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
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
