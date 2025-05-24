package syntax

import (
	"context"
	"regexp"
)

// ListNode represents a list block (ordered or unordered)
type ListNode struct {
	Items []*ListStructNode
}

type ListStructNode struct {
	Name  string // "ul" or "ol"
	Items []*ListItemNode
}

type ListItemNode struct {
	Children []interface{} // string or *ListStructNode
}

func (l *ListNode) ToHTML(ctx context.Context, inline Inline, options CallerOptions) string {
	html := ""
	for _, list := range l.Items {
		html += listToHTML(list, ctx, inline)
	}
	return html
}

func listToHTML(list *ListStructNode, ctx context.Context, inline Inline) string {
	html := "\n<" + list.Name + ">\n"
	for _, item := range list.Items {
		html += "<li>"
		for _, child := range item.Children {
			switch v := child.(type) {
			case string:
				html += inline.Format(v)
			case *ListStructNode:
				html += listToHTML(v, ctx, inline)
			}
		}
		html += "</li>\n"
	}
	html += "</" + list.Name + ">\n"
	return html
}

func (l *ListNode) AddChild(n Node)    {}
func (l *ListNode) GetContent() []Node { return nil }

type ListParser struct{}

var (
	reList = regexp.MustCompile(`^([-+]+)\s*(.+)`)
)

func (p *ListParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	line := scanner.Peek()
	if !reList.MatchString(line) {
		return false
	}
	var lines [][]string
	for !scanner.EOF() {
		l := scanner.Peek()
		if m := reList.FindStringSubmatch(l); m != nil {
			lines = append(lines, []string{l, m[1], m[2]})
			scanner.Next()
		} else {
			// リスト記法でない行は消費せずbreak（他パーサーに渡す）
			break
		}
	}
	if len(lines) == 0 {
		return false
	}
	var ret []*ListStructNode
	var listStack []*ListStructNode
	for _, row := range lines {
		// row[0]: line, row[1]: marks, row[2]: text
		marks, text := row[1], row[2]
		level := len(marks)
		// ul/ol 判定: 記号列の最後が + なら ol, それ以外は ul
		var typ string
		if marks[len(marks)-1] == '+' {
			typ = "ol"
		} else {
			typ = "ul"
		}
		if text == "" {
			continue
		}
		// スタックを調整
		for level < len(listStack) {
			listStack = listStack[:len(listStack)-1]
		}
		if level == len(listStack) && len(listStack) > 0 && listStack[len(listStack)-1].Name != typ {
			listStack = listStack[:len(listStack)-1]
		}
		for len(listStack) < level {
			container := &ListStructNode{Name: typ, Items: []*ListItemNode{}}
			if len(listStack) > 0 {
				lastItems := listStack[len(listStack)-1].Items
				if len(lastItems) > 0 {
					lastLi := lastItems[len(lastItems)-1]
					lastLi.Children = append(lastLi.Children, container)
				} else {
					item := &ListItemNode{Children: []interface{}{container}}
					listStack[len(listStack)-1].Items = append(listStack[len(listStack)-1].Items, item)
				}
			} else {
				ret = append(ret, container)
			}
			listStack = append(listStack, container)
		}
		item := &ListItemNode{Children: []interface{}{text}}
		listStack[len(listStack)-1].Items = append(listStack[len(listStack)-1].Items, item)
	}
	node := &ListNode{Items: ret}
	parent.AddChild(node)
	return true
}
