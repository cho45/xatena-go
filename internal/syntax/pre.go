package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
)

var PreTemplate = htmltpl.Must(htmltpl.New("pre").Parse(`
<pre>{{.Content}}</pre>
`))

// PreNode represents a <pre> block (stopp block with <pre> wrapper)
type PreNode struct {
	Content []Node // StopPNodeのように子ノードを持つ
}

func (p *PreNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	content := ContentToHTML(p, ctx, xatena, CallerOptions{
		stopp: true,
	})
	params := map[string]interface{}{"Content": htmltpl.HTML(content)}
	html := xatena.ExecuteTemplate("pre", params)
	return html
}

func (p *PreNode) AddChild(n Node) {
	p.Content = append(p.Content, n)
}

func (p *PreNode) GetContent() []Node {
	return p.Content
}

type PreParser struct{}

var rePreStart = regexp.MustCompile(`^>\|$`)
var rePreEnd = regexp.MustCompile(`^(.*?)\|<$`)

func (p *PreParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if scanner.Scan(rePreStart) {
		node := &PreNode{}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}
	if scanner.Scan(rePreEnd) {
		m := scanner.Matched()
		parent.AddChild(&TextNode{Text: m[1]})
		if len(*stack) == 0 {
			return false
		}
		*stack = (*stack)[:len(*stack)-1]
		return true
	}
	return false
}
