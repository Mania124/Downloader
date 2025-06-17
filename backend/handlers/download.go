package handlers

import (
	"context"
	"downloader/utils"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type DownloadRequest struct {
	URL        string `json:"url"`
	Format     string `json:"format"`
	Resolution string `json:"resolution"`
}

func DownloadVideo(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil || !utils.IsValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request or URL"})
		return
	}

	if req.Format != "video" && req.Format != "audio" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Choose 'video' or 'audio'"})
		return
	}

	outputPath := filepath.Join(utils.GetDownloadFolder(), "%(title)s.%(ext)s")
	var args []string

	if req.Format == "audio" {
		args = []string{
			"-f", "bestaudio/best",
			"--extract-audio", "--audio-format", "mp3", "--audio-quality", "192K",
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath, req.URL,
		}
	} else {
		format := "bestvideo+bestaudio/best"
		if req.Resolution != "" {
			format = fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best", req.Resolution)
		}
		args = []string{
			"-f", format,
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath, req.URL,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("yt-dlp error: %v\nOutput: %s", err, output)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
}
