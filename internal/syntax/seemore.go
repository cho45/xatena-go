package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
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
	html := xatena.ExecuteTemplate("seemore", params)
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

func (p *SeeMoreParser) CanHandle(line string) bool {
	return strings.HasPrefix(line, "====")
}

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
