package utils

import (
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
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return "audio"
	default:
		return "unknown"
	}
}

// GetContentType returns the MIME type for a file based on its extension
func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/avi"
	case ".mkv":
		return "video/mkv"
	case ".mov":
		return "video/mov"
	case ".wmv":
		return "video/wmv"
	case ".flv":
		return "video/flv"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mp3"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
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
