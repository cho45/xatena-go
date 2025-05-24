package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
	"text/template"

	"github.com/cho45/xatena-go/internal/util"
)

var SuperPreTemplate = htmltpl.Must(htmltpl.New("superpre").Parse(`
<pre class="{{.Class}}">{{.RawText}}</pre>
`))

// SuperPreNode represents a <pre> block with HTML-escaped content (super pre)
type SuperPreNode struct {
	Lang    string // e.g. "perl", "python" (optional)
	RawText string // raw preformatted text (will be HTML-escaped)
}

func (s *SuperPreNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	className := "code"
	langClass := ""
	if s.Lang != "" {
		langClass = " lang-" + s.Lang
	}
	params := map[string]interface{}{
		"Class":   className + langClass,
		"RawText": htmltpl.HTML(util.EscapeHTML(s.RawText)),
	}
	html, err := xatena.ExecuteTemplate("superpre", params)
	if err != nil {
		return `<div class="xatena-template-error">template error: ` + template.HTMLEscapeString(err.Error()) + `</div>`
	}
	return html
}

func (s *SuperPreNode) AddChild(n Node) {
	// SuperPreNodeは子ノードを持たない
}

func (s *SuperPreNode) GetContent() []Node {
	return nil
}

type SuperPreParser struct{}

var reSuperPreStart = regexp.MustCompile(`^>\|([^|]*)\|$`)
var reSuperPreEnd = regexp.MustCompile(`^\|\|<$`)

func (p *SuperPreParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if scanner.Scan(reSuperPreStart) {
		lang := scanner.Matched()[1]
		lines := scanner.ScanUntil(reSuperPreEnd)
		lines = lines[:len(lines)-1] // remove last matched
		node := &SuperPreNode{
			Lang:    lang,
			RawText: strings.Join(lines, "\n"),
		}
		parent.AddChild(node)
		return true
	}
	return false
}
