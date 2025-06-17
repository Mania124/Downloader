package handlers

import (
	"bufio"
	"context"
	"downloader/utils"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func DownloadWithProgress(c *gin.Context) {
	url := c.Query("url")
	format := c.Query("format")

	if url == "" || !utils.IsValidURL(url) || (format != "video" && format != "audio") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid url/format"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var args []string
	outputPath := filepath.Join(utils.GetDownloadFolder(), "%(title)s.%(ext)s")

	if format == "audio" {
		args = []string{
			"-f", "bestaudio/best",
			"--extract-audio", "--audio-format", "mp3", "--audio-quality", "192K",
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath,
			"--progress-template", "download:%(progress._percent_str)s (%(progress.eta)s remaining)",
			url,
		}
	} else {
		args = []string{
			"-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best",
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath,
			"--progress-template", "download:%(progress._percent_str)s (%(progress.eta)s remaining)",
			url,
		}

	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start process"})
		return
	}

	if err := cmd.Start(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start download"})
		return
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if _, err := c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", line)); err == nil {
				c.Writer.Flush()
			}
		}
	}

	cmd.Wait()
	c.Writer.WriteString("event: done\ndata: completed\n\n")
	c.Writer.Flush()
}
