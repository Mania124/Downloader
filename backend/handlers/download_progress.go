package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func DownloadWithProgress(c *gin.Context) {
	url := c.Query("url")
	format := c.Query("format") // video or audio

	if url == "" || (format != "video" && format != "audio") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid url/format"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var args []string
	outputPath := "%(title)s.%(ext)s"

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
			"-f", "bestvideo+bestaudio/best",
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath,
			"--progress-template", "download:%(progress._percent_str)s (%(progress.eta)s remaining)",
			url,
		}
	}

	cmd := exec.Command("yt-dlp", args...)
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
			c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", line))
			c.Writer.Flush()
		}
	}

	cmd.Wait()
	c.Writer.WriteString("event: done\ndata: completed\n\n")
	c.Writer.Flush()
}
