package xatena

import (
	"html"
	"regexp"
	"strings"
	"testing"
)

func TestXatenaCLIEquivalent(t *testing.T) {
	formatter := NewInlineFormatter()
	input := "[http://example.com:title]"
	output := formatter.Format(input)
	expected := `<a href="http://example.com">http://example.com</a>`
	if !strings.Contains(output, expected) {
		t.Errorf("expected %q to contain %q", output, expected)
	}
}

func TestInlineFormatterCustomTitleHandler(t *testing.T) {
	formatter := NewInlineFormatter(func(f *InlineFormatter) {
		f.SetTitleHandler(func(uri string) string {
			return "カスタムタイトル"
		})
	})
	input := "[http://example.com:title]"
	output := formatter.Format(input)
	expected := `<a href="http://example.com">カスタムタイトル</a>`
	if !strings.Contains(output, expected) {
		t.Errorf("expected %q to contain %q", output, expected)
	}
}

func TestInlineFormatterTitleEquals(t *testing.T) {
	formatter := NewInlineFormatter()
	input := "[http://example.com:title=foobar]"
	output := formatter.Format(input)
	expected := `<a href="http://example.com">foobar</a>`
	if !strings.Contains(output, expected) {
		t.Errorf("expected %q to contain %q", output, expected)
	}
}

func TestInlineFormatterAddRule(t *testing.T) {
	formatter := NewInlineFormatter()
	formatter.AddRule(InlineRule{
		Pattern: regexp.MustCompile(`\[custom:(.+?)\]`),
		Handler: func(f *InlineFormatter, m []string) string {
			return "<span class='custom'>" + html.EscapeString(m[1]) + "</span>"
		},
	})
	input := "[custom:テスト]"
	output := formatter.Format(input)
	expected := `<span class='custom'>テスト</span>`
	if !strings.Contains(output, expected) {
		t.Errorf("expected %q to contain %q", output, expected)
	}
}
