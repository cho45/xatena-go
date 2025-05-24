package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
)

var BlockquoteTemplate = htmltpl.Must(htmltpl.New("blockquote").Parse(`
<blockquote{{if .Cite}} cite="{{.Cite}}"{{end}}>
{{.Content}}
{{if .Title}}<cite>{{.Title}}</cite>{{end}}
</blockquote>
`))

// BlockquoteNode represents a blockquote block.
type BlockquoteNode struct {
	Cite    string // cite URL (optional)
	Content []Node // nested block nodes
}

func (b *BlockquoteNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	citeText := b.Cite
	title := ""
	uri := ""
	if citeText != "" {
		if isURL(citeText) {
			if strings.Contains(citeText, ":title=") {
				parts := strings.SplitN(citeText, ":title=", 2)
				uri = parts[0]
				titleText := parts[1]
				title = `<a href="` + uri + `">` + htmltpl.HTMLEscapeString(titleText) + `</a>`
			} else if strings.Contains(citeText, ":title") {
				uri = strings.SplitN(citeText, ":title", 2)[0]
				title = `<a href="` + uri + `">Example Web Page</a>`
			} else {
				uri = citeText
				title = xatena.GetInline().Format(ctx, "["+citeText+"]")
			}
		} else {
			title = xatena.GetInline().Format(ctx, citeText)
		}
		re := regexp.MustCompile(`href="([^"]+)"`)
		if m := re.FindStringSubmatch(title); m != nil {
			uri = m[1]
		} else if isURL(citeText) {
			uri = citeText
		}
	}
	content := ContentToHTML(b, ctx, xatena, options)

	html, err := xatena.ExecuteTemplate("blockquote", map[string]interface{}{
		"Cite":    uri,
		"Title":   htmltpl.HTML(title),
		"Content": htmltpl.HTML(content),
	})
	if err != nil {
		return `<div class="xatena-template-error">template error: ` + htmltpl.HTMLEscapeString(err.Error()) + `</div>`
	}
	return html
}

func (b *BlockquoteNode) GetContent() []Node {
	return b.Content
}

func (b *BlockquoteNode) AddChild(n Node) {
	b.Content = append(b.Content, n)
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

type BlockquoteParser struct{}

var reBlockquote = regexp.MustCompile(`^>([^>]+)?>$`)
var reBlockquoteEnd = regexp.MustCompile(`^<<$`)

func (p *BlockquoteParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	// BEGINNING: ^>(.*?)>$
	if scanner.Scan(reBlockquote) {
		node := &BlockquoteNode{}
		m := scanner.Matched()
		if len(m) > 1 {
			txt := strings.TrimSpace(m[1])
			node.Cite = txt
		}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}
	// ENDOFNODE: ^<<$
	if scanner.Scan(reBlockquoteEnd) {
		// Sectionノードを飛ばしてpop
		for len(*stack) > 0 {
			if _, ok := (*stack)[len(*stack)-1].(*SectionNode); ok {
				*stack = (*stack)[:len(*stack)-1]
			} else {
				break
			}
		}
		if len(*stack) == 0 {
			return false
		}
		*stack = (*stack)[:len(*stack)-1]
		return true
	}
	return false
}
