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
	if !scanner.Scan(reSection) {
		return false
	}
	m := scanner.Matched()
	stars := m[1]
	title := strings.TrimSpace(m[2])
	level := len(stars)
	if m[2] == "" {
		title = stars[1:]
		level = 1
	}
	if level > 3 {
		title = strings.Repeat("*", level-1) + title
		level = 1
	}
	sec := &SectionNode{Level: level, Title: title}
	for len(*stack) > 0 {
		if s, ok := (*stack)[len(*stack)-1].(*SectionNode); ok && s.Level >= level {
			*stack = (*stack)[:len(*stack)-1]
		} else if s, ok := (*stack)[len(*stack)-1].(*SeeMoreNode); ok && level == 1 && !s.IsSuper {
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

func (s *SectionNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	tmpl := `
<div class="section">
<h{{.Level}}>{{.Title}}</h{{.Level}}>
{{.Content}}
</div>`
	inline := xatena.GetInline()
	title := inline.Format(ctx, s.Title)
	content := ContentToHTML(s, ctx, xatena, options)
	var sb strings.Builder
	t := htmltpl.Must(htmltpl.New("section").Parse(tmpl))
	_ = t.Execute(&sb, map[string]interface{}{
		"Level":   s.Level + 2,
		"Title":   htmltpl.HTML(title),
		"Content": htmltpl.HTML(content),
	})
	return sb.String()
}
