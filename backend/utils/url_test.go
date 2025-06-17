package utils

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"ftp://example.com", true},
		{"example.com", false},         // missing scheme
		{"http:/example.com", false},   // invalid scheme format
		{"", false},                    // empty string
		{"ht!tp://example.com", false}, // invalid characters
		{"http://", false},             // missing host
		{"http://?query=123", false},   // missing host but has query
	}

	for _, test := range tests {
		result := IsValidURL(test.input)
		if result != test.expected {
			t.Errorf("IsValidURL(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
