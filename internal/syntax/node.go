package syntax

import (
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

type BlockNode interface {
	ToHTML() string
	AddChild(n BlockNode)
}

type RootNode struct {
	Content []BlockNode
}

func (r *RootNode) ToHTMLParagraph(text string) string {
	// Xatena.pm as_html_paragraph 互換: \n\n+ で分割し <p>...</p> で囲む、空行数に応じて <br /> を出力
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

func (r *RootNode) ToHTML() string {
	html := ""
	var textBuf []string
	flushParagraph := func() {
		if len(textBuf) > 0 {
			html += r.ToHTMLParagraph(strings.Join(textBuf, "\n"))
			textBuf = nil
		}
	}
	for _, n := range r.Content {
		if t, ok := n.(*TextNode); ok {
			textBuf = append(textBuf, t.Text)
		} else {
			flushParagraph()
			html += n.ToHTML()
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
	Parse(scanner *LineScanner, parent BlockNode, stack *[]BlockNode) bool
}

func (r *RootNode) AddChild(n BlockNode) {
	r.Content = append(r.Content, n)
}

type TextNode struct {
	Text string
}

func (t *TextNode) ToHTML() string {
	return util.EscapeHTML(t.Text)
}

func (t *TextNode) AddChild(n BlockNode) {
	panic("TextNode cannot have children")
}
