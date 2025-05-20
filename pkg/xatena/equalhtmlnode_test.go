package xatena

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestEqualHTMLNode(t *testing.T) {
	tests := []struct {
		a, b  string
		equal bool
		name  string
	}{
		{
			a:     `<div><p>foo</p></div>`,
			b:     `<div><p>foo</p></div>`,
			equal: true,
			name:  "identical simple HTML",
		},
		{
			a:     `<div><p>foo</p></div>`,
			b:     `<div> <p>foo</p> </div>`,
			equal: true,
			name:  "whitespace difference",
		},
		{
			a:     `<div><p>foo</p></div>`,
			b:     `<div><p>bar</p></div>`,
			equal: false,
			name:  "different text content",
		},
		{
			a:     `<div><p>foo</p></div>`,
			b:     `<div><span>foo</span></div>`,
			equal: false,
			name:  "different tag",
		},
		{
			a:     `<div class='a'><p>foo</p></div>`,
			b:     `<div class='a'><p>foo</p></div>`,
			equal: true,
			name:  "same attribute",
		},
		{
			a:     `<div class='a'><p>foo</p></div>`,
			b:     `<div class='b'><p>foo</p></div>`,
			equal: false,
			name:  "different attribute",
		},
	}

	for _, tt := range tests {
		nodeA, errA := parseHTMLNode(tt.a)
		nodeB, errB := parseHTMLNode(tt.b)
		if errA != nil || errB != nil {
			t.Fatalf("parse error: %v / %v", errA, errB)
		}
		result := equalHTMLNode(nodeA, nodeB)
		if result != tt.equal {
			t.Errorf("%s: expected %v, got %v\nA: %s\nB: %s", tt.name, tt.equal, result, tt.a, tt.b)
		}
	}
}

// parseHTMLNode is a helper for testing
func parseHTMLNode(s string) (*html.Node, error) {
	return html.Parse(strings.NewReader(s))
}
