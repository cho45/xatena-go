package syntax

import (
	"reflect"
	"testing"
)

func TestSplitForBreak(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single line no newline",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "single newline",
			input:    "hello\nworld",
			expected: []string{"hello", "world"},
		},
		{
			name:     "multiple newlines",
			input:    "hello\nworld\ntest",
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "trailing newline",
			input:    "hello\nworld\n",
			expected: []string{"hello", "world"},
		},
		{
			name:     "multiple trailing newlines",
			input:    "hello\nworld\n\n",
			expected: []string{"hello", "world"},
		},
		{
			name:     "leading newline",
			input:    "\nhello\nworld",
			expected: []string{"", "hello", "world"},
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: []string{},
		},
		{
			name:     "empty lines in middle",
			input:    "hello\n\nworld",
			expected: []string{"hello", "", "world"},
		},
		{
			name:     "single newline only",
			input:    "\n",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitForBreak(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SplitForBreak(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}