package syntax

import (
	"context"
	"regexp"
)

// ListNode represents a list block (ordered or unordered)
type ListNode struct {
	Ordered bool
	Items   []ListItemNode
}

type ListItemNode struct {
	Items []ListItemNode
}

func (l *ListNode) ToHTML(ctx context.Context) string {
	return ""
}

func (l *ListNode) AddChild(n Node) {
}

func (l *ListNode) GetContent() []Node {
	return nil
}

type ListParser struct{}

var (
	reListUnordered = regexp.MustCompile(`^(-+)(\s*)(.*)$`)
	reListOrdered   = regexp.MustCompile(`^(\++?)(\s*)(.*)$`)
)

func (p *ListParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	if !(reListUnordered.MatchString(line) || reListOrdered.MatchString(line)) {
		return false
	}
	// TODO
	return false
}
