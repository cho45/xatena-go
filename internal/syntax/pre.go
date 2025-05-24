package syntax

import (
	"context"
	"regexp"
)

// PreNode represents a <pre> block (stopp block with <pre> wrapper)
type PreNode struct {
	Content []Node // StopPNodeのように子ノードを持つ
}

func (p *PreNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	content := ContentToHTML(p, ctx, inline, CallerOptions{
		stopp: true,
	})
	return "<pre>" + content + "</pre>"
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
	line := scanner.Peek()
	if rePreStart.MatchString(line) {
		scanner.Next() // consume start
		node := &PreNode{}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}
	if m := rePreEnd.FindStringSubmatch(line); m != nil {
		parent.AddChild(&TextNode{Text: m[1]})
		scanner.Next() // consume end
		if len(*stack) == 0 {
			return false
		}
		*stack = (*stack)[:len(*stack)-1]
		return true
	}
	return false
}
