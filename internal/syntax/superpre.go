package syntax

import (
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

// SuperPreNode represents a <pre> block with HTML-escaped content (super pre)
type SuperPreNode struct {
	Lang    string // e.g. "perl", "python" (optional)
	RawText string // raw preformatted text (will be HTML-escaped)
}

func (s *SuperPreNode) ToHTML() string {
	class := ""
	if s.Lang != "" {
		class = " class=\"code lang-" + s.Lang + "\""
	}
	// HTMLエスケープ
	return "<pre" + class + ">\n" + util.EscapeHTML(s.RawText) + "\n</pre>"
}

func (s *SuperPreNode) AddChild(n BlockNode) {
	panic("SuperPreNode does not support adding child nodes")
}

type SuperPreParser struct{}

var reSuperPreStart = regexp.MustCompile(`^>\|\|([^|]*)\|?$`)
var reSuperPreEnd = regexp.MustCompile(`^\|\|<$`)

func (p *SuperPreParser) Parse(scanner *LineScanner, parent BlockNode, stack *[]BlockNode) bool {
	line := scanner.Peek()
	m := reSuperPreStart.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	scanner.Next() // consume start
	lang := strings.TrimSpace(m[1])
	var content []string
	for !scanner.EOF() {
		l := scanner.Peek()
		if reSuperPreEnd.MatchString(l) {
			scanner.Next() // consume end
			break
		}
		content = append(content, scanner.Next())
	}
	node := &SuperPreNode{Lang: lang, RawText: strings.Join(content, "\n")}
	if add, ok := parent.(interface{ AddChild(BlockNode) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}
