package util

import "testing"

func TestCleanAddr(t *testing.T) {
	tests := []struct {
		name          string
		ssl           bool
		url           string
		expected      string
		errorExpected bool
	}{
		{
			name:     "without ssl",
			ssl:      false,
			url:      "example.com",
			expected: "http://example.com",
		},
		{
			name:     "with ssl",
			ssl:      true,
			url:      "example.com",
			expected: "https://example.com",
		},
		{
			name:     "with additional parts",
			ssl:      true,
			url:      "example.com/path/to/source",
			expected: "https://example.com",
		},
		{
			name:     "with schema",
			ssl:      true,
			url:      "http://example.com/path/to/source",
			expected: "http://example.com",
		},
		{
			name:          "empty url",
			ssl:           false,
			url:           "",
			errorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := CleanAddr(tt.url, tt.ssl)

			if err != nil {
				if !tt.errorExpected {
					t.Fatalf("expected nil error, got %q", err.Error())
				}
			} else {
				if addr != tt.expected {
					t.Fatalf("expected %q, got %q", tt.expected, addr)
				}
			}
		})
	}
}
