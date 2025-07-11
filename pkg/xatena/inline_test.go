package xatena

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

func TestInlineFormatter_Format(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "simple link",
			input:  "[http://example.com/]",
			expect: `<a href="http://example.com/">http://example.com/</a>`,
		},
		{
			name:   "simple link raw",
			input:  "http://example.com/",
			expect: `<a href="http://example.com/">http://example.com/</a>`,
		},
		{
			name:   "mailto",
			input:  "[mailto:foo@example.com]",
			expect: `<a href="mailto:foo@example.com">foo@example.com</a>`,
		},
		{
			name:   "tex",
			input:  "[tex:E=mc^2]",
			expect: `<img src="http://chart.apis.google.com/chart?cht=tx&chl=E%3Dmc%5E2" alt="E=mc^2"/>`,
		},
		{
			name:   "footnote",
			input:  "((note))",
			expect: `<a href="#fn1" title="note">*1</a>`,
		},
		{
			name:   "barcode",
			input:  "[http://example.com/:barcode]",
			expect: `<img src="http://chart.apis.google.com/chart?chs=150x150&cht=qr&chl=http%3A%2F%2Fexample.com%2F" title="http://example.com/"/>`,
		},
		{
			name:   "title option",
			input:  "[http://example.com/:title=Example]",
			expect: `<a href="http://example.com/">Example</a>`,
		},
		{
			name:   "unlink",
			input:  "[]unlinked[]",
			expect: "unlinked",
		},
	}

	for _, c := range cases {
		f := NewInlineFormatter()
		got := f.Format(context.Background(), c.input)
		if got != c.expect {
			t.Errorf("%s: input=%q\nexpect=%q\ngot   =%q", c.name, c.input, c.expect, got)
		}
	}
}

// TestInlineFormatterAddRuleAt tests AddRuleAt functionality
func TestInlineFormatterAddRuleAt(t *testing.T) {
	f := NewInlineFormatter()
	
	// Define a custom rule
	customRule := InlineRule{
		Pattern: regexp.MustCompile(`\[test:(.+?)\]`),
		Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
			return "<test>" + m[1] + "</test>"
		},
	}
	
	// Test adding at the beginning (index 0)
	f.AddRuleAt(0, customRule)
	result := f.Format(context.Background(), "[test:content]")
	expected := "<test>content</test>"
	if result != expected {
		t.Errorf("AddRuleAt(0): expected %q, got %q", expected, result)
	}
	
	// Test adding at invalid negative index (should append)
	f2 := NewInlineFormatter()
	f2.AddRuleAt(-1, customRule)
	result2 := f2.Format(context.Background(), "[test:content2]")
	expected2 := "<test>content2</test>"
	if result2 != expected2 {
		t.Errorf("AddRuleAt(-1): expected %q, got %q", expected2, result2)
	}
	
	// Test adding at index beyond length (should append)
	f3 := NewInlineFormatter()
	f3.AddRuleAt(100, customRule)
	result3 := f3.Format(context.Background(), "[test:content3]")
	expected3 := "<test>content3</test>"
	if result3 != expected3 {
		t.Errorf("AddRuleAt(100): expected %q, got %q", expected3, result3)
	}
}

// TestInlineFormatterFootnotes tests Footnotes functionality
func TestInlineFormatterFootnotes(t *testing.T) {
	f := NewInlineFormatter()
	
	// Test with input containing footnotes
	input := "テキスト((脚注内容))とさらに((別の脚注))"
	output := f.Format(context.Background(), input)
	
	// Get footnotes
	footnotes := f.Footnotes()
	
	// We expect 2 footnotes
	if len(footnotes) != 2 {
		t.Errorf("expected 2 footnotes, got %d", len(footnotes))
	}
	
	// Verify footnotes are captured correctly
	if len(footnotes) >= 1 {
		if footnotes[0].Note != "脚注内容" {
			t.Errorf("expected first footnote note to be '脚注内容', got %q", footnotes[0].Note)
		}
	}
	
	if len(footnotes) >= 2 {
		if footnotes[1].Note != "別の脚注" {
			t.Errorf("expected second footnote note to be '別の脚注', got %q", footnotes[1].Note)
		}
	}
	
	// Test that footnotes are included in output
	if !strings.Contains(output, "脚注内容") && !strings.Contains(output, "別の脚注") {
		t.Errorf("expected output to contain footnote content, got %q", output)
	}
}

// TestInlineFormatterFootnotesEmpty tests Footnotes with no footnotes
func TestInlineFormatterFootnotesEmpty(t *testing.T) {
	f := NewInlineFormatter()
	
	// Test with input containing no footnotes
	input := "普通のテキスト"
	f.Format(context.Background(), input)
	
	// Get footnotes
	footnotes := f.Footnotes()
	
	// We expect 0 footnotes
	if len(footnotes) != 0 {
		t.Errorf("expected 0 footnotes, got %d", len(footnotes))
	}
}
