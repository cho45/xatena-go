package syntax

import (
	"context"
	"regexp"
	"strings"
)

type CallerOptions struct {
	stopp bool
}

type Node interface {
	ToHTML(ctx context.Context, inline Inline, options CallerOptions) string
}

type HasContent interface {
	AddChild(n Node)
	GetContent() []Node
}

type RootNode struct {
	Content []Node
}

func SplitForBreak(text string) []string {
	if text == "" {
		return []string{}
	}
	parts := strings.Split(text, "\n")
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func ToHTMLParagraph(text string, inline Inline, options CallerOptions) string {
	text = inline.Format(text)
	if options.stopp {
		return text
	}
	re := regexp.MustCompile(`(\n{2,})`)
	parts := reSplitWithSep(re, text)

	html := "<p>"
	for _, para := range parts {
		if regexp.MustCompile(`^\n+$`).MatchString(para) {
			html += "</p>" + strings.Repeat("<br />\n", (len(para)-2)) + "<p>"
		} else {
			html += strings.Join(SplitForBreak(para), "<br />\n")
		}
	}
	html += "</p>"
	return html
}

func ContentToHTML(r HasContent, ctx context.Context, inline Inline, options CallerOptions) string {
	html := ""
	var textBuf []string
	flushParagraph := func() {
		hasText := strings.Join(textBuf, "") != ""
		if hasText {
			html += ToHTMLParagraph(strings.Join(textBuf, "\n"), inline, options)
		}
		textBuf = nil
	}
	for _, n := range r.GetContent() {
		if t, ok := n.(*TextNode); ok {
			textBuf = append(textBuf, t.Text)
		} else {
			flushParagraph()
			html += n.ToHTML(ctx, inline, options)
		}
	}
	flushParagraph()
	return html
}

// reSplitWithSep splits s by the regexp re, including the separator (match) in the result slice.
// e.g. reSplitWithSep(regexp.MustCompile(`(\n{2,})`), "a\n\nb\n\n\nc") => ["a", "\n\n", "b", "\n\n\n", "c"]
func reSplitWithSep(re *regexp.Regexp, s string) []string {
	result := []string{}
	last := 0
	indices := re.FindAllStringIndex(s, -1)
	for _, idx := range indices {
		if last < idx[0] {
			result = append(result, s[last:idx[0]])
		}
		result = append(result, s[idx[0]:idx[1]])
		last = idx[1]
	}
	if last < len(s) {
		result = append(result, s[last:])
	}
	return result
}

type BlockParser interface {
	Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool
}

func (r *RootNode) AddChild(n Node) {
	r.Content = append(r.Content, n)
}

func (r *RootNode) GetContent() []Node {
	return r.Content
}

func (r *RootNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	return ContentToHTML(r, ctx, inline, options)
}

type TextNode struct {
	Text string
}

func (t *TextNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	panic("TextNode does not support ToHTML")
}
