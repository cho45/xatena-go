package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
)

var reSection = regexp.MustCompile(`^(\*+)(\s*.*)$`)

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
	stars := m[1]
	title := strings.TrimSpace(m[2])
	level := len(stars)
	// スペースがない場合（"**", "***" など）は level=1, title="*"または"**" などにする
	if m[2] == "" {
		title = stars[1:]
		level = 1
	}
	if level > 3 {
		// ****foo のような場合も level=1, title=***foo
		title = strings.Repeat("*", level-1) + title
		level = 1
	}
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
	*stack = append(*stack, sec)
	return true
}

func (s *SectionNode) ToHTML(ctx context.Context, inline Inline) string {
	tmpl := `
<div class="section">
<h{{.Level}}>{{.Title}}</h{{.Level}}>
{{.Content}}
</div>`
	title := inline.Format(s.Title)
	content := ContentToHTML(s, ctx, inline)
	var sb strings.Builder
	t := htmltpl.Must(htmltpl.New("section").Parse(tmpl))
	_ = t.Execute(&sb, map[string]interface{}{
		"Level":   s.Level + 2,
		"Title":   htmltpl.HTML(title),
		"Content": htmltpl.HTML(content),
	})
	return sb.String()
}
