package handlers

import (
	"downloader/utils"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DownloadRequest struct {
	URL        string `json:"url"`
	Format     string `json:"format"`     // "video" or "audio"
	Resolution string `json:"resolution"` // optional
}

func DownloadVideo(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request or missing URL"})
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

	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(output)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Download completed", "log": string(output)})
}
