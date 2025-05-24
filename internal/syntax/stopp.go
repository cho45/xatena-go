package syntax

import (
	"context"
	"regexp"
)

// StopPNode represents a block that disables auto <p>/<br> insertion.
type StopPNode struct {
	Content []Node
}

func (s *StopPNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	return ContentToHTML(s, ctx, xatena, CallerOptions{
		stopp: true,
	})
}
func (s *StopPNode) AddChild(n Node) {
	s.Content = append(s.Content, n)
}

func (s *StopPNode) GetContent() []Node {
	return s.Content
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
