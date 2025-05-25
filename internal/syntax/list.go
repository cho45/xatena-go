package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
)

var ListTemplate = htmltpl.Must(htmltpl.New("list").Parse(`
{{.OpenTag}}
{{range .Items}}
  <li>{{range .Content}}{{.}}{{end}}</li>
{{end}}
{{.CloseTag}}
`))

// ListNode represents a list block (ordered or unordered)
type ListNode struct {
	Items []*ListStructNode
}

type ListStructNode struct {
	Name  string // "ul" or "ol"
	Items []*ListItemNode
}

type ListItemNode struct {
	Content []interface{} // string or *ListStructNode
}

func (l *ListNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	var html string
	for _, list := range l.Items {
		html += listStructToHTML(list, ctx, xatena)
	}
	return html
}

// 1つのListStructNode（ul/ol）を再帰的にHTML化
func listStructToHTML(list *ListStructNode, ctx context.Context, xatena XatenaContext) string {
	var items []map[string]interface{}
	for _, item := range list.Items {
		var content []interface{}
		for _, child := range item.Content {
			switch v := child.(type) {
			case string:
				content = append(content, htmltpl.HTML(xatena.GetInline().Format(ctx, v)))
			case *ListStructNode:
				content = append(content, htmltpl.HTML(listStructToHTML(v, ctx, xatena)))
			}
		}
		items = append(items, map[string]interface{}{"Content": content})
	}
	html := xatena.ExecuteTemplate("list", map[string]interface{}{
		"OpenTag":  htmltpl.HTML("<" + list.Name + ">"),
		"CloseTag": htmltpl.HTML("</" + list.Name + ">"),
		"Items":    items,
	})
	return html
}

func (l *ListNode) AddChild(n Node)    {}
func (l *ListNode) GetContent() []Node { return nil }

type ListParser struct{}

var (
	reList = regexp.MustCompile(`^([-+]+)\s*(.+)`)
)

func (p *ListParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if !scanner.Scan(reList) {
		return false
	}
	var lines [][]string
	m := scanner.Matched()
	lines = append(lines, []string{m[0], m[1], m[2]})
	for !scanner.EOF() {
		if !scanner.Scan(reList) {
			break
		}
		m := scanner.Matched()
		lines = append(lines, []string{m[0], m[1], m[2]})
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
					lastLi.Content = append(lastLi.Content, container)
				} else {
					item := &ListItemNode{Content: []interface{}{container}}
					listStack[len(listStack)-1].Items = append(listStack[len(listStack)-1].Items, item)
				}
			} else {
				ret = append(ret, container)
			}
			listStack = append(listStack, container)
		}
		item := &ListItemNode{Content: []interface{}{text}}
		listStack[len(listStack)-1].Items = append(listStack[len(listStack)-1].Items, item)
	}
	node := &ListNode{Items: ret}
	parent.AddChild(node)
	return true
}
