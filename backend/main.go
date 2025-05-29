package main

import (
	"downloader/router"
	"downloader/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Setup download folder if needed
	downloadFolder := utils.GetDownloadFolder()
	if _, err := os.Stat(downloadFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadFolder, os.ModePerm); err != nil {
			log.Fatalf("Failed to create download folder: %v", err)
		}
	}

	// Register routes
	router.SetupRoutes(r)

	r.Run(":5000")
}
