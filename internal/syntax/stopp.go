package syntax

import (
	"context"
	"regexp"
	"strings"
)

// StopPNode represents a block that disables auto <p>/<br> insertion.
type StopPNode struct {
	Content string // raw content (HTML allowed)
}

func (s *StopPNode) ToHTML(ctx context.Context, inline Inline) string {
	return s.Content
}
func (s *StopPNode) AddChild(n Node) {
	panic("StopPNode does not support child nodes")
}

func (s *StopPNode) GetContent() []Node {
	return nil
}

type StopPParser struct{}

var reStopPStart = regexp.MustCompile(`^>\s*$`)
var reStopPEnd = regexp.MustCompile(`^<\s*$`)

func (p *StopPParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	if !reStopPStart.MatchString(line) {
		return false
	}
	scanner.Next() // consume start
	var content []string
	for !scanner.EOF() {
		l := scanner.Peek()
		if reStopPEnd.MatchString(l) {
			scanner.Next() // consume end
			break
		}
		content = append(content, scanner.Next())
	}
	node := &StopPNode{Content: strings.Join(content, "\n")}
	if add, ok := parent.(interface{ AddChild(Node) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}
