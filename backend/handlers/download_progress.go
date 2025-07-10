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
	resolution := c.Query("resolution")
	videoFormat := c.Query("videoFormat")

	if url == "" || !utils.IsValidURL(url) || (format != "video" && format != "audio") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid url/format"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Get list of files before download to identify new files
	downloadFolder := utils.GetDownloadFolder()
	beforeFiles, _ := utils.GetFileList(downloadFolder)

	var args []string
	outputPath := filepath.Join(downloadFolder, "%(title)s.%(ext)s")

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
		format := utils.BuildVideoFormat(resolution, videoFormat)
		args = []string{
			"-f", format,
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath,
			"--progress-template", "download:%(progress._percent_str)s (%(progress.eta)s remaining)",
			url,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
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

	// Get list of files after download to identify the new file
	afterFiles, err := utils.GetFileList(downloadFolder)
	if err == nil {
		// Find the newly downloaded file
		newFiles := utils.FindNewFiles(beforeFiles, afterFiles)
		if len(newFiles) > 0 {
			downloadedFile := newFiles[0]
			c.Writer.WriteString(fmt.Sprintf("event: file\ndata: {\"filename\":\"%s\",\"downloadUrl\":\"/files/%s\"}\n\n", downloadedFile, downloadedFile))
			c.Writer.Flush()
		}
	}

	c.Writer.WriteString("event: done\ndata: completed\n\n")
	c.Writer.Flush()
}
