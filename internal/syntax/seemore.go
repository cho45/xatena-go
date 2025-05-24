package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
)

var SeeMoreTemplate = htmltpl.Must(htmltpl.New("seemore").Parse(`
<div class="seemore">{{.Content}}</div>
`))

// SeeMoreNode represents a <div class="seemore"> block
// (==== or ===== line)
type SeeMoreNode struct {
	IsSuper bool
	Content []Node
}

func (s *SeeMoreNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	content := ContentToHTML(s, ctx, xatena, options)
	params := map[string]interface{}{"Content": htmltpl.HTML(content)}
	html, err := xatena.ExecuteTemplate("seemore", params)
	if err != nil {
		return `<div class="xatena-template-error">template error: ` + htmltpl.HTMLEscapeString(err.Error()) + `</div>`
	}
	return html
}

func (s *SeeMoreNode) AddChild(n Node) {
	s.Content = append(s.Content, n)
}

func (s *SeeMoreNode) GetContent() []Node {
	return s.Content
}

type SeeMoreParser struct{}

var reSeeMore = regexp.MustCompile(`^====(=)?$`)

func (p *SeeMoreParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if scanner.Scan(reSeeMore) {
		isSuper := scanner.Matched()[1] != ""
		node := &SeeMoreNode{IsSuper: isSuper}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}

	return false
}
