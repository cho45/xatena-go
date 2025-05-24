package syntax

import (
	"context"
	"regexp"
)

type CommentNode struct{}

func (c *CommentNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	return "<!-- -->"
}
func (c *CommentNode) AddChild(n Node)    {}
func (c *CommentNode) GetContent() []Node { return nil }

type CommentParser struct{}

func (p *CommentParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	reBegin := regexp.MustCompile(`^(.*)<!--.*?(-->)?$`)
	reEnd := regexp.MustCompile(`^-->$`)
	if scanner.Scan(reBegin) {
		m := scanner.Matched()
		pre := m[1]
		if pre != "" {
			parent.AddChild(&TextNode{Text: pre})
		}
		if m[2] == "" {
			scanner.ScanUntil(reEnd)
		}
		parent.AddChild(&CommentNode{})
		return true
	}
	return false
}
