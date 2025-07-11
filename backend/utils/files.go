package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo represents information about a downloaded file
type FileInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ModTime     string `json:"modTime"`
	DownloadURL string `json:"downloadUrl"`
	Type        string `json:"type"`
}

// GetFileList returns a list of filenames in the given directory
func GetFileList(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// FindNewFiles returns files that are in 'after' but not in 'before'
func FindNewFiles(before, after []string) []string {
	beforeSet := make(map[string]bool)
	for _, file := range before {
		beforeSet[file] = true
	}

	var newFiles []string
	for _, file := range after {
		if !beforeSet[file] {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}

// GetFileType determines the type of file based on its extension
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".m4v", ".3gp", ".ogv":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma":
		return "audio"
	default:
		return "unknown"
	}
}

// GetContentType returns the MIME type for a file based on its extension
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".m4v":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".webm":
		return "video/webm"
	case ".3gp":
		return "video/3gpp"
	case ".ogv":
		return "video/ogg"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".aac", ".m4a":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	case ".wma":
		return "audio/x-ms-wma"
	default:
		return "application/octet-stream"
	}
}

// CreateFileInfo creates a FileInfo struct for a given file
func CreateFileInfo(filename string, downloadFolder string) (*FileInfo, error) {
	filePath := filepath.Join(downloadFolder, filename)
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:        filename,
		Size:        info.Size(),
		ModTime:     info.ModTime().Format(time.RFC3339),
		DownloadURL: "/files/" + filename,
		Type:        GetFileType(filename),
	}, nil
}

// CleanupOldFiles removes files older than the specified duration
func CleanupOldFiles(downloadFolder string, maxAge time.Duration) error {
	entries, err := os.ReadDir(downloadFolder)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(downloadFolder, entry.Name())
			if err := os.Remove(filePath); err != nil {
				// Log error but continue with other files
				continue
			}
		}
	}
	return nil
}

// BuildVideoFormat constructs the yt-dlp format string based on resolution and video format preferences
func BuildVideoFormat(resolution, videoFormat string) string {
	// Handle video format preference with more specific selectors
	var formatOptions []string

	// Handle resolution preference
	var heightFilter string
	if resolution != "" {
		heightFilter = fmt.Sprintf("[height<=%s]", resolution)
	}

	// Build format string based on preferences
	switch videoFormat {
	case "mp4":
		if resolution != "" {
			formatOptions = []string{
				fmt.Sprintf("bestvideo[ext=mp4]%s+bestaudio[ext=m4a]/bestvideo[ext=mp4]%s+bestaudio/best[ext=mp4]%s", heightFilter, heightFilter, heightFilter),
				fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
			}
		} else {
			formatOptions = []string{
				"bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo[ext=mp4]+bestaudio/best[ext=mp4]",
				"bestvideo+bestaudio/best",
			}
		}
	case "webm":
		if resolution != "" {
			formatOptions = []string{
				fmt.Sprintf("bestvideo[ext=webm]%s+bestaudio[ext=webm]/best[ext=webm]%s", heightFilter, heightFilter),
				fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
			}
		} else {
			formatOptions = []string{
				"bestvideo[ext=webm]+bestaudio[ext=webm]/best[ext=webm]",
				"bestvideo+bestaudio/best",
			}
		}
	case "mkv":
		if resolution != "" {
			formatOptions = []string{
				fmt.Sprintf("bestvideo[ext=mkv]%s+bestaudio/best[ext=mkv]%s", heightFilter, heightFilter),
				fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
			}
		} else {
			formatOptions = []string{
				"bestvideo[ext=mkv]+bestaudio/best[ext=mkv]",
				"bestvideo+bestaudio/best",
			}
		}
	default:
		// For other formats or "best", use general approach
		if resolution != "" {
			formatOptions = []string{
				fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
			}
		} else {
			formatOptions = []string{
				"bestvideo+bestaudio/best",
			}
		}
	}

	// Join all options with fallbacks
	return strings.Join(formatOptions, "/")
}
