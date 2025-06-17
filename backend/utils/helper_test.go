package utils

import (
	"os"
	"testing"
)

func TestGetDownloadFolder(t *testing.T) {
	// Save original env
	origUserProfile := os.Getenv("USERPROFILE")
	origHome := os.Getenv("HOME")

	defer func() {
		os.Setenv("USERPROFILE", origUserProfile)
		os.Setenv("HOME", origHome)
	}()

	os.Setenv("USERPROFILE", "/test/userprofile")
	os.Setenv("HOME", "/test/home")

	tests := []struct {
		goos     string
		expected string
	}{
		{"windows", "/test/userprofile/Downloads"},
		{"darwin", "/test/home/Downloads"},
		{"linux", "/test/home/Downloads"},
	}

	for _, test := range tests {
		t.Run(test.goos, func(t *testing.T) {
			getRuntimeGOOS = func() string { return test.goos }

			result := GetDownloadFolder()
			if result != test.expected {
				t.Errorf("GetDownloadFolder() for %s = %s; want %s", test.goos, result, test.expected)
			}
		})
	}
}
