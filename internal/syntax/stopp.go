package syntax

import (
	"context"
	"regexp"
)

// StopPNode represents a block that disables auto <p>/<br> insertion.
type StopPNode struct {
	children []Node
}

func (s *StopPNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	return ContentToHTML(s, ctx, inline, CallerOptions{
		stopp: true,
	})
}
func (s *StopPNode) AddChild(n Node) {
	s.children = append(s.children, n)
}

func (s *StopPNode) GetContent() []Node {
	return s.children
}

type StopPParser struct{}

var reStopPStart = regexp.MustCompile(`^>(<.+>)(<)?$`)
var reStopPEnd = regexp.MustCompile(`^(.+>)<`)

func (p *StopPParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if scanner.Scan(reStopPStart) {
		node := &StopPNode{}
		node.AddChild(&TextNode{Text: scanner.Matched()[1]}) // Add the opening tag
		parent.AddChild(node)
		if scanner.Matched()[2] == "" {
			*stack = append(*stack, node)
		}
		return true
	}

	if scanner.Scan(reStopPEnd) {
		lastParent := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		lastParent.AddChild(&TextNode{Text: scanner.Matched()[1]})
		return true
	}
	return false
}
