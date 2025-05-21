package syntax

import (
	"context"
	"regexp"
	"strings"
)

// PreNode represents a <pre> block (with optional language class)
type PreNode struct {
	Lang    string // e.g. "perl", "python" (optional)
	RawText string // raw preformatted text
}

func (p *PreNode) ToHTML(ctx context.Context, inline Inline) string {
	class := ""
	if p.Lang != "" {
		class = " class=\"code lang-" + p.Lang + "\""
	}
	return "<pre" + class + ">\n" + p.RawText + "\n</pre>"
}

func (p *PreNode) AddChild(n Node) {
	panic("PreNode does not support adding child nodes")
}

func (p *PreNode) GetContent() []Node {
	return nil
}

type PreParser struct{}

var rePreStart = regexp.MustCompile(`^>\|([^|]*)\|?$`)
var rePreEnd = regexp.MustCompile(`^\|<$`)

func (p *PreParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	m := rePreStart.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	scanner.Next() // consume start
	lang := strings.TrimSpace(m[1])
	var content []string
	for !scanner.EOF() {
		l := scanner.Peek()
		if rePreEnd.MatchString(l) {
			scanner.Next() // consume end
			break
		}
		content = append(content, scanner.Next())
	}
	node := &PreNode{Lang: lang, RawText: strings.Join(content, "\n")}
	if add, ok := parent.(interface{ AddChild(Node) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}
