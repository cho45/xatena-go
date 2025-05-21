package syntax

import (
	"context"
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

var reDefinitionList = regexp.MustCompile(`^:([^:]+):(.*)$`)
var reDefinitionListCont = regexp.MustCompile(`^::(.*)$`)

// DefinitionListNode represents a definition list block
type DefinitionListNode struct {
	Items []DefinitionItemNode
}

type DefinitionItemNode struct {
	Term string
	Desc string
}

func (d *DefinitionListNode) ToHTML(ctx context.Context, inline Inline) string {
	html := "<dl>\n"
	for _, item := range d.Items {
		html += "  <dt>" + util.EscapeHTML(item.Term) + "</dt>\n"
		html += "  <dd>" + util.EscapeHTML(item.Desc) + "</dd>\n"
	}
	html += "</dl>"
	return html
}

func (d *DefinitionListNode) AddChild(n Node) {
	// 定義リストは子ブロックを持たないので何もしない
}

// DefinitionListNode（仮名）などContentを持つ型の場合:
func (d *DefinitionListNode) GetContent() []Node {
	return nil
}

// DefinitionListParser parses definition list blocks
// 例: :term:desc で始まる連続行をまとめてパース
type DefinitionListParser struct{}

func (p *DefinitionListParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if !reDefinitionList.MatchString(scanner.Peek()) {
		return false
	}
	var items []DefinitionItemNode
	for !scanner.EOF() {
		line := scanner.Peek()
		m := reDefinitionList.FindStringSubmatch(line)
		if m == nil {
			break
		}
		scanner.Next()
		term := strings.TrimSpace(m[1])
		desc := strings.TrimSpace(m[2])
		// handle multiline desc
		for !scanner.EOF() {
			m2 := reDefinitionListCont.FindStringSubmatch(scanner.Peek())
			if m2 == nil {
				break
			}
			scanner.Next()
			desc += "\n" + strings.TrimSpace(m2[1])
		}
		items = append(items, DefinitionItemNode{Term: term, Desc: desc})
	}
	if len(items) == 0 {
		return false
	}

	node := &DefinitionListNode{Items: items}
	if add, ok := parent.(interface{ AddChild(Node) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}
