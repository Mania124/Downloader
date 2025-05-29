package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type ThumbnailRequest struct {
	URL string `json:"url"`
}

func GetThumbnail(c *gin.Context) {
	var req ThumbnailRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request or missing URL"})
		return
	}

	cmd := exec.Command("yt-dlp", "-J", req.URL)
	output, err := cmd.Output()
	if err != nil {
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
