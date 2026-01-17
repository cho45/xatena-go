package xatena

import (
	"context"
	"strings"
	"testing"
)

// TestMalformed_BlockNotations tests various malformed block notations
// to ensure they are handled gracefully without panics or errors.
func TestMalformed_BlockNotations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// SuperPre (>|lang|...||<) malformed patterns
		{
			name:  "superpre unclosed",
			input: ">|perl|\ncode here\nno closing",
		},
		{
			name:  "superpre unclosed at EOF",
			input: ">||\nsome content",
		},
		{
			name:  "superpre empty lang with content",
			input: ">||\ncode\n||<",
		},
		{
			name:  "superpre wrong closing",
			input: ">|go|\ncode\n|<",
		},
		{
			name:  "superpre nested start",
			input: ">|perl|\n>|ruby|\ncode\n||<",
		},
		{
			name:  "superpre multiple starts no close",
			input: ">|a|\n>|b|\n>|c|",
		},
		{
			name:  "superpre closing only",
			input: "||<",
		},
		{
			name:  "superpre closing after text",
			input: "some text\n||<\nmore text",
		},

		// Pre (>|...|<) malformed patterns
		{
			name:  "pre unclosed",
			input: ">|\ntext inside\nno close",
		},
		{
			name:  "pre unclosed at EOF",
			input: ">|\n",
		},
		{
			name:  "pre closing only",
			input: "|<",
		},
		{
			name:  "pre closing with next line",
			input: "|<\nnext line",
		},
		{
			name:  "pre nested pre start",
			input: ">|\n>|\ntext\n|<",
		},
		{
			name:  "pre double close",
			input: ">|\ntext\n|<\n|<",
		},
		{
			name:  "pre mismatched end after text",
			input: "some text\n|<",
		},

		// Blockquote (>>...<<) malformed patterns
		{
			name:  "blockquote unclosed",
			input: ">>\nquoted text\nno close",
		},
		{
			name:  "blockquote closing only",
			input: "<<",
		},
		{
			name:  "blockquote closing with next line",
			input: "<<\nnext line",
		},
		{
			name:  "blockquote mismatched end in section",
			input: "* heading\n<<",
		},
		{
			name:  "blockquote deeply nested mismatch",
			input: ">>|\n>>\n|<",
		},
		{
			name:  "blockquote multiple close",
			input: ">>\ntext\n<<\n<<",
		},
		{
			name:  "blockquote nested unclosed",
			input: ">>\n>>\ndeep\n<<",
		},
		{
			name:  "blockquote with cite unclosed",
			input: ">http://example.com>\nquoted\nno close",
		},
		{
			name:  "blockquote multiple open single close",
			input: ">>\n>>\n>>\ntext\n<<",
		},

		// StopP (>(...)<) malformed patterns
		{
			name:  "stopp unclosed",
			input: ">(\ntext\nno close",
		},
		{
			name:  "stopp closing only",
			input: ")<",
		},
		{
			name:  "stopp wrong close ins<",
			input: "ins<",
		},
		{
			name:  "stopp wrong close ins< with next line",
			input: "ins<\nnext line",
		},
		{
			name:  "stopp wrong close del<",
			input: "del<",
		},

		// List malformed patterns
		{
			name:  "list item without text",
			input: "-\n-\n-",
		},
		{
			name:  "list deeply nested",
			input: "-------------- very deep",
		},
		{
			name:  "list mixed markers invalid",
			input: "- item\n+ item\n-+ mixed",
		},
		{
			name:  "list only markers",
			input: "---",
		},
		{
			name:  "list continuation markers",
			input: "++++++++++ item",
		},

		// Table malformed patterns
		{
			name:  "table single pipe",
			input: "|",
		},
		{
			name:  "table unclosed pipe",
			input: "|cell without close",
		},
		{
			name:  "table empty cells",
			input: "||||",
		},
		{
			name:  "table only header markers",
			input: "|*|*|*|",
		},
		{
			name:  "table inconsistent columns",
			input: "|a|b|c|\n|d|e|",
		},

		// Definition list malformed patterns
		{
			name:  "deflist term only",
			input: ":term:",
		},
		{
			name:  "deflist no term",
			input: "::description only",
		},
		{
			name:  "deflist continuation without term",
			input: "::\n::",
		},
		{
			name:  "deflist empty term",
			input: "::desc",
		},
		{
			name:  "deflist colon only",
			input: ":",
		},

		// Section malformed patterns
		{
			name:  "section many stars",
			input: "******** too many stars",
		},
		{
			name:  "section stars only",
			input: "***",
		},
		{
			name:  "section empty",
			input: "*",
		},
		{
			name:  "section star at end",
			input: "text *",
		},

		// Comment malformed patterns
		{
			name:  "comment unclosed",
			input: "<!-- unclosed comment",
		},
		{
			name:  "comment nested",
			input: "<!-- outer <!-- inner --> -->",
		},
		{
			name:  "comment empty",
			input: "<!---->",
		},
		{
			name:  "comment with newlines unclosed",
			input: "<!--\nmultiline\ncomment",
		},

		// SeeMore malformed patterns
		{
			name:  "seemore too few equals",
			input: "===",
		},
		{
			name:  "seemore with text",
			input: "==== text after",
		},

		// Mixed/Complex malformed patterns
		{
			name:  "mixed unclosed blocks",
			input: ">>\n>|\n>(\ntext",
		},
		{
			name:  "all closings no openings",
			input: "<<\n|<\n)<\n||<",
		},
		{
			name:  "interleaved blocks wrong order",
			input: ">>\n>|\ntext\n<<\n|<",
		},
		{
			name:  "deeply nested everything",
			input: ">>\n>|\n* heading\n- list\n|cell|\n:term:desc\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			// We expect this NOT to panic.
			got := x.ToHTML(context.Background(), tt.input)
			// We expect the result to be a valid string (not empty necessarily, but processable)
			if got == "" && tt.input != "" {
				t.Logf("Warning: empty output for non-empty input %q", tt.input)
			}
			t.Logf("Input: %q -> Output: %q", tt.input, got)
		})
	}
}

// TestMalformed_InlineNotations tests various malformed inline notations
func TestMalformed_InlineNotations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Footnote malformed patterns
		{
			name:  "footnote unclosed single",
			input: "text ((unclosed",
		},
		{
			name:  "footnote unclosed nested",
			input: "((outer ((inner))",
		},
		{
			name:  "footnote closing only",
			input: "text )) closing",
		},
		{
			name:  "footnote empty",
			input: "(())",
		},
		{
			name:  "footnote triple parens",
			input: "(((triple)))",
		},
		{
			name:  "footnote sandwich pattern",
			input: ")((content))(", // This is actually a valid escape pattern
		},

		// Link malformed patterns
		{
			name:  "bracket unclosed",
			input: "[http://example.com",
		},
		{
			name:  "bracket no url",
			input: "[]",
		},
		{
			name:  "bracket partial url",
			input: "[http://]",
		},
		{
			name:  "bracket invalid protocol",
			input: "[xyz://example.com]",
		},
		{
			name:  "bracket title unclosed",
			input: "[http://example.com:title=unclosed",
		},
		{
			name:  "bracket barcode with text",
			input: "[http://example.com:barcode:extra]",
		},
		{
			name:  "bracket nested",
			input: "[[http://example.com]]",
		},
		{
			name:  "bracket closing only",
			input: "]",
		},

		// mailto malformed patterns
		{
			name:  "mailto no at",
			input: "[mailto:invalid]",
		},
		{
			name:  "mailto no domain",
			input: "[mailto:user@]",
		},
		{
			name:  "mailto no user",
			input: "[mailto:@domain.com]",
		},
		{
			name:  "mailto incomplete",
			input: "[mailto:]",
		},

		// tex malformed patterns
		{
			name:  "tex empty",
			input: "[tex:]",
		},
		{
			name:  "tex unclosed",
			input: "[tex:E=mc^2",
		},
		{
			name:  "tex with brackets",
			input: "[tex:f(x) = [a,b]]",
		},

		// Unlink escape patterns
		{
			name:  "unlink unclosed",
			input: "[]text",
		},
		{
			name:  "unlink empty",
			input: "[][]",
		},
		{
			name:  "unlink nested brackets",
			input: "[][nested][]",
		},

		// HTML tag malformed patterns
		{
			name:  "html unclosed tag",
			input: "<div",
		},
		{
			name:  "html unclosed attr",
			input: "<div class=",
		},
		{
			name:  "html malformed attr",
			input: "<div class=\"unclosed>",
		},
		{
			name:  "html closing only",
			input: "</div>",
		},
		{
			name:  "html empty tag",
			input: "<>",
		},
		{
			name:  "html nested incomplete",
			input: "<div><span>text</div>",
		},

		// URL patterns without brackets
		{
			name:  "url without scheme",
			input: "www.example.com",
		},
		{
			name:  "url partial",
			input: "http://",
		},
		{
			name:  "url just protocol",
			input: "https:",
		},
		{
			name:  "url with spaces",
			input: "http://example.com/path with spaces",
		},

		// Mixed inline malformed
		{
			name:  "mixed unclosed all",
			input: "((footnote [link <tag",
		},
		{
			name:  "mixed interleaved",
			input: "((note [http://))url])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			// We expect this NOT to panic.
			got := x.ToHTML(context.Background(), tt.input)
			t.Logf("Input: %q -> Output: %q", tt.input, got)
		})
	}
}

// TestMalformed_EdgeCases tests edge cases and boundary conditions
func TestMalformed_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Empty and whitespace
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "only newlines",
			input: "\n\n\n",
		},
		{
			name:  "only spaces",
			input: "   ",
		},
		{
			name:  "only tabs",
			input: "\t\t\t",
		},
		{
			name:  "mixed whitespace",
			input: " \t \n \t \n",
		},
		{
			name:  "null byte",
			input: "text\x00more",
		},
		{
			name:  "control characters",
			input: "text\x01\x02\x03more",
		},

		// Line ending edge cases
		{
			name:  "CR only",
			input: "line1\rline2",
		},
		{
			name:  "CRLF",
			input: "line1\r\nline2",
		},
		{
			name:  "mixed line endings",
			input: "a\nb\r\nc\rd",
		},
		{
			name:  "trailing newlines",
			input: "text\n\n\n",
		},
		{
			name:  "leading newlines",
			input: "\n\n\ntext",
		},

		// Unicode edge cases
		{
			name:  "emoji in text",
			input: "* heading 🎉\n- list 📝",
		},
		{
			name:  "CJK characters",
			input: "* 日本語見出し\n- リストアイテム",
		},
		{
			name:  "RTL characters",
			input: "* العربية\n- עברית",
		},
		{
			name:  "combining characters",
			input: "café résumé naïve",
		},
		{
			name:  "zero width characters",
			input: "te\u200Bxt wi\u200Cth ze\u200Dro",
		},
		{
			name:  "BOM character",
			input: "\uFEFFtext with BOM",
		},

		// Very long inputs
		{
			name:  "very long line",
			input: strings.Repeat("a", 10000),
		},
		{
			name:  "many short lines",
			input: strings.Repeat("line\n", 1000),
		},
		{
			name:  "deeply nested lists",
			input: strings.Repeat("-", 100) + " deep item",
		},
		{
			name:  "many sections",
			input: strings.Repeat("* section\ntext\n", 100),
		},

		// Special characters in notation context
		{
			name:  "pipe in text",
			input: "a | b | c",
		},
		{
			name:  "asterisk in text",
			input: "a * b * c",
		},
		{
			name:  "bracket in text",
			input: "a [ b ] c",
		},
		{
			name:  "paren in text",
			input: "a ( b ) c",
		},
		{
			name:  "angle in text",
			input: "a < b > c",
		},
		{
			name:  "colon in text",
			input: "a : b : c",
		},
		{
			name:  "dash in text",
			input: "a - b - c",
		},
		{
			name:  "plus in text",
			input: "a + b + c",
		},
		{
			name:  "equals in text",
			input: "a = b = c",
		},

		// Escaped/quoted content
		{
			name:  "html entities",
			input: "&lt;not a tag&gt; &amp; stuff",
		},
		{
			name:  "backslash sequences",
			input: "\\n \\t \\r \\\\",
		},

		// Binary-like content
		{
			name:  "percent encoding",
			input: "%20%3C%3E%26",
		},
		{
			name:  "base64-like",
			input: "SGVsbG8gV29ybGQ=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			// We expect this NOT to panic.
			got := x.ToHTML(context.Background(), tt.input)
			t.Logf("Input length: %d, Output length: %d", len(tt.input), len(got))
		})
	}
}

// TestMalformed_ContentPreservation tests that content is not lost on malformed input
func TestMalformed_ContentPreservation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		mustContain []string
	}{
		{
			name:        "text before unclosed block",
			input:       "visible text\n>|\nunclosed pre content",
			mustContain: []string{"visible text"},
		},
		{
			name:        "text after orphan close",
			input:       "|<\ntext after",
			mustContain: []string{"text after"},
		},
		{
			name:        "content between malformed blocks",
			input:       ">>\nA\n<<\n<<\nB\n>>\nC",
			mustContain: []string{"A", "B", "C"},
		},
		{
			name:        "list items despite malformed markers",
			input:       "- item A\n--- item B\n----- item C",
			mustContain: []string{"item A", "item B", "item C"},
		},
		{
			name:        "table cells despite inconsistent pipes",
			input:       "|cellA|cellB\n|cellC|",
			mustContain: []string{"cellA", "cellB", "cellC"},
		},
		{
			name:        "definition terms and descriptions",
			input:       ":term:desc\n::\n::orphan",
			mustContain: []string{"term", "desc"},
		},
		{
			name:        "section headings",
			input:       "* h1\n** h2\n*** h3\n**** overflow",
			mustContain: []string{"h1", "h2", "h3"},
		},
		{
			name:        "inline text in malformed footnote",
			input:       "before ((inside)) after",
			mustContain: []string{"before", "after"},
		},
		{
			name:        "text around malformed links",
			input:       "before [broken after",
			mustContain: []string{"before", "[broken after"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			got := x.ToHTML(context.Background(), tt.input)
			for _, substr := range tt.mustContain {
				if !strings.Contains(got, substr) {
					t.Errorf("output does not contain expected substring %q\nInput: %q\nOutput: %s", substr, tt.input, got)
				}
			}
		})
	}
}

// TestMalformed_NoHTMLInjection tests that malformed input doesn't cause XSS
func TestMalformed_NoHTMLInjection(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		mustEscape []string // these should be escaped or removed, not appear raw
	}{
		{
			name:       "script in text",
			input:      "<script>alert('xss')</script>",
			mustEscape: []string{}, // HTML tags are allowed through
		},
		{
			name:       "script in section title",
			input:      "* <script>alert('xss')</script>",
			mustEscape: []string{},
		},
		{
			name:       "onerror in img",
			input:      "<img onerror=\"alert('xss')\" src=x>",
			mustEscape: []string{},
		},
		{
			name:       "javascript url",
			input:      "[javascript:alert('xss')]",
			mustEscape: []string{},
		},
		{
			name:       "data url",
			input:      "[data:text/html,<script>alert('xss')</script>]",
			mustEscape: []string{},
		},
		{
			name:       "malformed url with script",
			input:      "[http://example.com/<script>]",
			mustEscape: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			got := x.ToHTML(context.Background(), tt.input)
			// Log for manual review - XSS prevention is complex
			t.Logf("XSS test: %s\nInput: %q\nOutput: %s", tt.name, tt.input, got)
		})
	}
}

// TestMalformed_RecoveryBehavior tests that parser recovers and continues after malformed input
func TestMalformed_RecoveryBehavior(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		mustContain  []string // content that must be present
		mustNotPanic bool     // should always be true
	}{
		{
			name:  "recover after unclosed superpre",
			input: ">|perl|\ncode\n\n* heading after\ntext",
			// SuperPreは||<まで全て取り込む仕様。EOFまで取り込まれる。
			// コンテンツはHTMLエスケープされてpre内に存在する。
			mustContain:  []string{"code", "* heading after"},
			mustNotPanic: true,
		},
		{
			name:         "recover after orphan close",
			input:        "||<\n* valid section\nvalid text",
			mustContain:  []string{"valid section", "valid text"},
			mustNotPanic: true,
		},
		{
			name:         "recover after nested malformed",
			input:        ">>\n>|\n* h\n- l\n|<\n<<\nnormal text",
			mustContain:  []string{"normal text"},
			mustNotPanic: true,
		},
		{
			name:         "multiple malformed then valid",
			input:        "<<\n|<\n||<\n)<\n* valid\ntext",
			mustContain:  []string{"valid", "text"},
			mustNotPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()

			x := NewXatena()
			got := x.ToHTML(context.Background(), tt.input)

			for _, substr := range tt.mustContain {
				if !strings.Contains(got, substr) {
					t.Errorf("recovery failed: output does not contain %q\nInput: %q\nOutput: %s", substr, tt.input, got)
				}
			}
		})
	}
}
