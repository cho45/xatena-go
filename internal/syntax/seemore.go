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
	content := ContentToHTML(s, ctx, inline, options)
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
	if scanner.Scan(reSeeMore) {
		isSuper := scanner.Matched()[1] != ""
		node := &SeeMoreNode{IsSuper: isSuper}
		parent.AddChild(node)
		*stack = append(*stack, node)
		return true
	}

	return false
}
