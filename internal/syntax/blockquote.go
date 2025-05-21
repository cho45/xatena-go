package syntax

import (
	"context"
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
	for _, n := range b.Content {
		html += n.ToHTML(ctx)
	}
	if b.Title != "" {
		// Xatena仕様: cite内容がURLならリンク化
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

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

type BlockquoteParser struct{}

var reBlockquote = regexp.MustCompile(`^>([^>]+)?>$`)

func (p *BlockquoteParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	return false
}

type BlockquoteLineNode struct{}

func (b *BlockquoteLineNode) GetContent() []Node {
	return nil
}
