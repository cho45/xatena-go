package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
)

var CommentTemplate = htmltpl.Must(htmltpl.New("comment").Parse(`
{{.Content}}
`))

type CommentNode struct{}

func (c *CommentNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	html := xatena.ExecuteTemplate("comment", map[string]interface{}{
		"Content": htmltpl.HTML("<!-- -->"),
	})
	return html
}
func (c *CommentNode) AddChild(n Node)    {}
func (c *CommentNode) GetContent() []Node { return nil }

type CommentParser struct{}

var reBegin = regexp.MustCompile(`^(.*)<!--(.*)$`)

func (p *CommentParser) CanHandle(line string) bool {
	return true
}

func (p *CommentParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if scanner.Scan(reBegin) {
		m := scanner.Matched()
		pre := m[1]
		rest := m[2]

		if pre != "" {
			parent.AddChild(&TextNode{Text: pre})
		}

		// Consume until "-->"
		if strings.Contains(rest, "-->") {
			// Ends on the same line
		} else {
			// Scan until "-->" or EOF
			for !scanner.EOF() {
				line := scanner.Next()
				if strings.Contains(line, "-->") {
					break
				}
			}
		}

		parent.AddChild(&CommentNode{})
		return true
	}
	return false
}
