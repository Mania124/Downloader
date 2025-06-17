package handlers

import (
	"context"
	"downloader/utils"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

type ThumbnailRequest struct {
	URL string `json:"url"`
}

var YTDLPCommand = func(ctx context.Context, url string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "yt-dlp", "-J", url)
	return cmd.Output()
}

func GetThumbnail(c *gin.Context) {
	var req ThumbnailRequest
	if err := c.ShouldBindJSON(&req); err != nil || !utils.IsValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	output, err := YTDLPCommand(ctx, req.URL)
	if err != nil {
		log.Printf("yt-dlp error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video info"})
		return
	}

	var metadata struct {
		Thumbnail string `json:"thumbnail"`
	}
	if err := json.Unmarshal(output, &metadata); err != nil || metadata.Thumbnail == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Thumbnail not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"thumbnail": metadata.Thumbnail})
}
