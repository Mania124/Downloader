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
	// Default format preferences
	var formatSelectors []string

	// Handle video format preference
	var extFilter string
	switch videoFormat {
	case "mp4":
		extFilter = "[ext=mp4]"
	case "webm":
		extFilter = "[ext=webm]"
	case "mkv":
		extFilter = "[ext=mkv]"
	case "avi":
		extFilter = "[ext=avi]"
	case "mov":
		extFilter = "[ext=mov]"
	case "flv":
		extFilter = "[ext=flv]"
	case "3gp":
		extFilter = "[ext=3gp]"
	default:
		extFilter = "" // No extension filter, let yt-dlp choose best
	}

	// Handle resolution preference
	var heightFilter string
	if resolution != "" {
		heightFilter = fmt.Sprintf("[height<=%s]", resolution)
	}

	// Build format string with preferences
	if resolution != "" && videoFormat != "" && videoFormat != "best" {
		// Specific resolution and format - prioritize exact match
		formatSelectors = append(formatSelectors,
			fmt.Sprintf("bestvideo%s%s+bestaudio/best%s%s", heightFilter, extFilter, heightFilter, extFilter),
			fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
		)
	} else if resolution != "" {
		// Specific resolution, any format
		formatSelectors = append(formatSelectors,
			fmt.Sprintf("bestvideo%s+bestaudio/best%s", heightFilter, heightFilter),
		)
	} else if videoFormat != "" && videoFormat != "best" {
		// Specific format, any resolution - prioritize format
		formatSelectors = append(formatSelectors,
			fmt.Sprintf("bestvideo%s+bestaudio/best%s", extFilter, extFilter),
		)
	} else {
		// Default: best quality available
		formatSelectors = append(formatSelectors, "bestvideo+bestaudio/best")
	}

	// Return the first (most preferred) format selector, or default if none
	if len(formatSelectors) > 0 {
		return formatSelectors[0]
	}
	return "bestvideo+bestaudio/best"
}
