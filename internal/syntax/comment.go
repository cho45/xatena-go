package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
)

var CommentTemplate = htmltpl.Must(htmltpl.New("comment").Parse(`
{{.Content}}
`))

type CommentNode struct{}

func (c *CommentNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	html, err := xatena.ExecuteTemplate("comment", map[string]interface{}{
		"Content": htmltpl.HTML("<!-- -->"),
	})
	if err != nil {
		return `<div class="xatena-template-error">template error: ` + htmltpl.HTMLEscapeString(err.Error()) + `</div>`
	}
	return html
}
func (c *CommentNode) AddChild(n Node)    {}
func (c *CommentNode) GetContent() []Node { return nil }

type CommentParser struct{}

func (p *CommentParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	reBegin := regexp.MustCompile(`^(.*)<!--.*?(-->)?$`)
	reEnd := regexp.MustCompile(`^-->$`)
	if scanner.Scan(reBegin) {
		m := scanner.Matched()
		pre := m[1]
		if pre != "" {
			parent.AddChild(&TextNode{Text: pre})
		}
		if m[2] == "" {
			scanner.ScanUntil(reEnd)
		}
		parent.AddChild(&CommentNode{})
		return true
	}
	return false
}
