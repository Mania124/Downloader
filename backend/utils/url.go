package utils

import "net/url"

func IsValidURL(rawURL string) bool {
	parsed, err := url.ParseRequestURI(rawURL)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}
