package syntax

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

var reSection = regexp.MustCompile(`^(\*+)(\s*)(.*)$`)

// SectionNode represents a section (heading + content)
type SectionNode struct {
	Level   int    // 1=*, 2=**, ...
	Title   string // heading text
	Content []Node // nested block nodes
}

func (s *SectionNode) AddChild(n Node) {
	s.Content = append(s.Content, n)
}

func (s *SectionNode) GetContent() []Node {
	return s.Content
}

type SectionTitleNode struct{}

func (s *SectionTitleNode) GetContent() []Node {
	return nil
}

type SectionParser struct{}

func (p *SectionParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	m := reSection.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	scanner.Next() // consume heading
	level := len(m[1])
	title := strings.TrimSpace(m[3])
	sec := &SectionNode{Level: level, Title: title}
	// stackを巻き戻して親を決定（Xatena.pm互換）
	for len(*stack) > 0 {
		if s, ok := (*stack)[len(*stack)-1].(*SectionNode); ok && s.Level >= level {
			*stack = (*stack)[:len(*stack)-1]
		} else {
			break
		}
	}
	parentNode := (*stack)[len(*stack)-1]
	parentNode.AddChild(sec)
	// 新しいセクションをstackに積む
	*stack = append(*stack, sec)
	return true
}

func (s *SectionNode) ToHTML(ctx context.Context) string {
	html := `<div class="section">\n` +
		fmt.Sprintf("<h%d>%s</h%d>\n", s.Level+2, util.EscapeHTML(s.Title), s.Level+2)
	html += "</div>"
	return html
}
