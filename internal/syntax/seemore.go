package syntax

import (
	"context"
	"regexp"
)

// SeeMoreNode represents a <div class="seemore"> block
// (==== or ===== line)
type SeeMoreNode struct {
	IsSuper  bool
	children []Node
}

func (s *SeeMoreNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	content := ""
	for _, child := range s.children {
		content += child.ToHTML(ctx, inline, options)
	}
	return `<div class="seemore">` + content + `</div>`
}

func (s *SeeMoreNode) AddChild(n Node) {
	s.children = append(s.children, n)
}

func (s *SeeMoreNode) GetContent() []Node {
	return s.children
}

type SeeMoreParser struct{}

var reSeeMore = regexp.MustCompile(`^====(=)?$`)

func (p *SeeMoreParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	m := reSeeMore.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	scanner.Next() // consume
	node := &SeeMoreNode{IsSuper: m[1] != ""}
	parent.AddChild(node)
	*stack = append(*stack, node)
	return true
}
