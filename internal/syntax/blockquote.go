package syntax

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

// BlockquoteNode represents a blockquote block.
type BlockquoteNode struct {
	Cite    string // cite URL (optional)
	Title   string // title or label (optional)
	Content []Node // nested block nodes
}

func (b *BlockquoteNode) ToHTML(ctx context.Context) string {
	html := "<blockquote"
	if b.Cite != "" {
		html += " cite=\"" + util.EscapeHTML(b.Cite) + "\""
	}
	html += ">\n"
	html += ContentToHTML(b, ctx)
	fmt.Printf("BlockquoteNode: %v\n", ContentToHTML(b, ctx))
	if b.Title != "" {
		if isURL(b.Title) {
			html += "<cite><a href=\"" + util.EscapeHTML(b.Title) + "\">" + util.EscapeHTML(b.Title) + "</a></cite>\n"
		} else {
			html += "<cite>" + util.EscapeHTML(b.Title) + "</cite>\n"
		}
	}
	html += "</blockquote>"
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

func (p *BlockquoteParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	// BEGINNING: ^>(.*?)>$
	if m := reBlockquote.FindStringSubmatch(scanner.Peek()); m != nil {
		scanner.Next()
		node := &BlockquoteNode{}
		if len(m) > 1 {
			txt := strings.TrimSpace(m[1])
			if strings.HasPrefix(txt, "http://") || strings.HasPrefix(txt, "https://") {
				node.Cite = txt
				node.Title = txt
			} else {
				node.Title = txt
			}
		}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}
	// ENDOFNODE: ^<<$
	if strings.TrimSpace(scanner.Peek()) == "<<" {
		scanner.Next()
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
