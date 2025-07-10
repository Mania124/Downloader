package handlers

import (
	"context"
	"downloader/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type DownloadRequest struct {
	URL         string `json:"url"`
	Format      string `json:"format"`      // "video" or "audio"
	Resolution  string `json:"resolution"`  // "360", "480", "720", "1080"
	VideoFormat string `json:"videoFormat"` // "mp4", "webm", "mkv", "avi", "best"
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

	// Get list of files before download to identify new files
	downloadFolder := utils.GetDownloadFolder()
	beforeFiles, err := utils.GetFileList(downloadFolder)
	if err != nil {
		log.Printf("Error getting file list before download: %v", err)
	}

	outputPath := filepath.Join(downloadFolder, "%(title)s.%(ext)s")
	var args []string

	if req.Format == "audio" {
		args = []string{
			"-f", "bestaudio/best",
			"--extract-audio", "--audio-format", "mp3", "--audio-quality", "192K",
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath, req.URL,
		}
	} else {
		format := utils.BuildVideoFormat(req.Resolution, req.VideoFormat)
		args = []string{
			"-f", format,
			"--no-playlist", "--prefer-free-formats",
			"-o", outputPath, req.URL,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("yt-dlp error: %v\nOutput: %s", err, output)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed"})
		return
	}

	// Get list of files after download to identify the new file
	afterFiles, err := utils.GetFileList(downloadFolder)
	if err != nil {
		log.Printf("Error getting file list after download: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		return
	}

	// Find the newly downloaded file
	newFiles := utils.FindNewFiles(beforeFiles, afterFiles)
	if len(newFiles) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		return
	}

	// Return information about the downloaded file
	downloadedFile := newFiles[0] // Take the first new file
	fileInfo, err := os.Stat(filepath.Join(downloadFolder, downloadedFile))
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Download completed",
		"filename":    downloadedFile,
		"size":        fileInfo.Size(),
		"downloadUrl": fmt.Sprintf("/files/%s", downloadedFile),
	})
}

// ServeFile serves a downloaded file to the client
func ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)
	filePath := filepath.Join(utils.GetDownloadFolder(), filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Get file info for content length
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Set appropriate headers
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Set content type based on file extension
	c.Header("Content-Type", utils.GetContentType(filename))

	// Stream the file to the client
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		log.Printf("Error streaming file %s: %v", filePath, err)
		return
	}
}

// ListFiles returns a list of all downloaded files
func ListFiles(c *gin.Context) {
	downloadFolder := utils.GetDownloadFolder()

	// Read directory contents
	entries, err := os.ReadDir(downloadFolder)
	if err != nil {
		log.Printf("Error reading download folder: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read download folder"})
		return
	}

	var files []utils.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		// Create file info using utility function
		fileInfo, err := utils.CreateFileInfo(entry.Name(), downloadFolder)
		if err != nil {
			log.Printf("Error getting file info for %s: %v", entry.Name(), err)
			continue
		}

		files = append(files, *fileInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"count": len(files),
	})
}
