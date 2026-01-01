package xatena

import (
	"context"
	"testing"
)

func TestPanic_IsolatedClosingNotations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "isolated pre end",
			input: "|<\nnext line",
		},
		{
			name:  "isolated blockquote end",
			input: "<<\nnext line",
		},
		{
			name:  "isolated stopp end",
			input: "ins<\nnext line",
		},
		{
			name:  "mismatched end after text",
			input: "some text\n|<",
		},
		{
			name:  "mismatched end in section",
			input: "* heading\n<<",
		},
		{
			name:  "deeply nested mismatch",
			input: ">>|\n>>\n|<", // |< closes the pre, but what happens to the nested blockquote?
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			// We expect this NOT to panic. If it does, the test fails.
			res := x.ToHTML(context.Background(), tt.input)
			t.Logf("Result for %s: %q", tt.name, res)
		})
	}
}
