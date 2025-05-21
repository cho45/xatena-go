package syntax

import (
	"context"
	"regexp"
	"strings"
)

type Node interface {
	ToHTML(ctx context.Context, inline Inline) string
}

type HasContent interface {
	AddChild(n Node)
	GetContent() []Node
}

type RootNode struct {
	Content []Node
}

func ToHTMLParagraph(text string, inline Inline) string {
	text = inline.Format(text)
	re := regexp.MustCompile(`(\n{2,})`)
	parts := reSplitWithSep(re, text)

	html := "<p>"
	for _, para := range parts {
		if regexp.MustCompile(`^\n+$`).MatchString(para) {
			html += "</p>" + strings.Repeat("<br />\n", (len(para)-2)) + "<p>"
		} else {
			html += strings.Join(strings.Split(para, "\n"), "<br />\n")
		}
	}
	html += "</p>"
	return html
}

func ContentToHTML(r HasContent, ctx context.Context, inline Inline) string {
	html := ""
	var textBuf []string
	flushParagraph := func() {
		if len(textBuf) > 0 {
			html += ToHTMLParagraph(strings.Join(textBuf, "\n"), inline)
			textBuf = nil
		}
	}
	for _, n := range r.GetContent() {
		if t, ok := n.(*TextNode); ok {
			textBuf = append(textBuf, t.Text)
		} else {
			flushParagraph()
			html += n.ToHTML(ctx, inline)
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

func (r *RootNode) ToHTML(ctx context.Context, inline Inline) string {
	return ContentToHTML(r, ctx, inline)
}

type TextNode struct {
	Text string
}

func (t *TextNode) ToHTML(ctx context.Context, inline Inline) string {
	panic("TextNode does not support ToHTML")
}
