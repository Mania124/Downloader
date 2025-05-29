package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetDownloadFolder() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	default:
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	}
}
