package util

import (
	"testing"
)

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic HTML tags",
			input:    "<div>Hello</div>",
			expected: "&lt;div&gt;Hello&lt;/div&gt;",
		},
		{
			name:     "quotes",
			input:    `"Hello World"`,
			expected: "&#34;Hello World&#34;",
		},
		{
			name:     "apostrophe",
			input:    "It's a test",
			expected: "It&#39;s a test",
		},
		{
			name:     "ampersand",
			input:    "Tom & Jerry",
			expected: "Tom &amp; Jerry",
		},
		{
			name:     "mixed special characters",
			input:    `<script>alert("XSS");</script>`,
			expected: "&lt;script&gt;alert(&#34;XSS&#34;);&lt;/script&gt;",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "unicode characters",
			input:    "こんにちは世界",
			expected: "こんにちは世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHTML(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
