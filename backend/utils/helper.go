package utils

import (
	"os"
	"path/filepath"
)

var getRuntimeGOOS = func() string { return os.Getenv("GOOS_OVERRIDE") }

func GetDownloadFolder() string {
	goos := getRuntimeGOOS()
	if goos == "" {
		goos = runtimeGOOS
	}

	switch goos {
	case "windows":
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	default:
		return filepath.Join(os.Getenv("HOME"), "Downloads")
	}
}

// This allows us to override GOOS during tests
var runtimeGOOS = detectGOOS()

func detectGOOS() string {
	return os.Getenv("GOOS_REAL") // unused during normal run, set only in tests
}
