package syntax

import (
	"context"
	"regexp"
)

var reDefinitionList = regexp.MustCompile(`^:([^:]+):(.*)$`)
var reDefinitionListCont = regexp.MustCompile(`^::(.*)$`)

// DefinitionListNode represents a definition list block
type DefinitionListNode struct {
	Items []DefinitionItemNode
}

type DefinitionItemNode struct {
	Term  string
	Descs []string // 複数のddを保持
}

func (d *DefinitionListNode) ToHTML(ctx context.Context, inline Inline) string {
	html := "<dl>\n"
	for _, item := range d.Items {
		html += "  <dt>" + inline.Format(item.Term) + "</dt>\n"
		for _, desc := range item.Descs {
			html += "  <dd>" + inline.Format(desc) + "</dd>\n"
		}
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
	var lines []string
	matched := false
	for !scanner.EOF() {
		line := scanner.Peek()
		if reDefinitionList.MatchString(line) || reDefinitionListCont.MatchString(line) {
			lines = append(lines, scanner.Next())
			matched = true
		} else {
			break
		}
	}
	if !matched {
		return false
	}

	items := []DefinitionItemNode{}
	var currentTerm string
	var currentDescs []string
	for _, line := range lines {
		if m := reDefinitionListCont.FindStringSubmatch(line); m != nil {
			// ::description
			if len(currentDescs) > 0 {
				currentDescs[len(currentDescs)-1] += m[1]
			} else {
				currentDescs = append(currentDescs, m[1])
			}
		} else if m := reDefinitionList.FindStringSubmatch(line); m != nil {
			// :term:desc
			if currentTerm != "" {
				items = append(items, DefinitionItemNode{Term: currentTerm, Descs: currentDescs})
			}
			currentTerm = m[1]
			currentDescs = []string{}
			desc := m[2]
			if desc != "" {
				currentDescs = append(currentDescs, desc)
			}
		}
	}
	if currentTerm != "" {
		items = append(items, DefinitionItemNode{Term: currentTerm, Descs: currentDescs})
	}
	node := &DefinitionListNode{Items: items}
	if add, ok := parent.(interface{ AddChild(Node) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}
