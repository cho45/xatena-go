package syntax

import (
	"context"
	"regexp"
	"strings"
)

// ListNode represents a list block (ordered or unordered)
type ListNode struct {
	Ordered bool
	Items   []ListItemNode
}

type ListItemNode struct {
	Content []BlockNode // each item can contain blocks (usually ParagraphNode)
}

func (l *ListNode) ToHTML(ctx context.Context) string {
	tag := "ul"
	if l.Ordered {
		tag = "ol"
	}
	html := "<" + tag + ">\n"
	for _, item := range l.Items {
		html += "  <li>"
		// TextNodeをまとめてパラグラフ化、ListNodeはそのままネスト
		var textBuf []string
		for _, n := range item.Content {
			if t, ok := n.(*TextNode); ok {
				textBuf = append(textBuf, t.Text)
			} else {
				if len(textBuf) > 0 {
					textBuf = nil
				}
				html += n.ToHTML(ctx)
			}
		}
		if len(textBuf) > 0 {
		}
		html += "</li>\n"
	}
	html += "</" + tag + ">"
	return html
}

func (l *ListNode) AddChild(n BlockNode) {
	if len(l.Items) == 0 {
		l.Items = append(l.Items, ListItemNode{})
	}
	last := &l.Items[len(l.Items)-1]
	last.Content = append(last.Content, n)
}

// ListParser parses list blocks (ordered/unordered)
type ListParser struct{}

var (
	reListUnordered = regexp.MustCompile(`^(-+)(\s*)(.*)$`)
	reListOrdered   = regexp.MustCompile(`^(\++?)(\s*)(.*)$`)
)

func (p *ListParser) Parse(scanner *LineScanner, parent BlockNode, stack *[]BlockNode) bool {
	line := scanner.Peek()
	if !(reListUnordered.MatchString(line) || reListOrdered.MatchString(line)) {
		return false
	}

	type listStackElem struct {
		level   int
		ordered bool
		list    *ListNode
		item    *ListItemNode
	}
	var lstack []*listStackElem
	var rootList *ListNode

	for !scanner.EOF() {
		l := scanner.Peek()
		var level int
		var ordered bool
		var content string
		if m := reListUnordered.FindStringSubmatch(l); m != nil {
			level = len(m[1])
			ordered = false
			content = m[3]
		} else if m := reListOrdered.FindStringSubmatch(l); m != nil {
			level = len(m[1])
			ordered = true
			content = m[3]
		} else {
			// interruption: not a list line, break and unwind
			break
		}
		scanner.Next()

		// Stack unwinding for level/type change
		for len(lstack) > 0 && (lstack[len(lstack)-1].level > level || lstack[len(lstack)-1].ordered != ordered) {
			lstack = lstack[:len(lstack)-1]
		}
		// Stack push for deeper level or new type
		if len(lstack) == 0 || lstack[len(lstack)-1].level < level || lstack[len(lstack)-1].ordered != ordered {
			ln := &ListNode{Ordered: ordered}
			var item *ListItemNode
			if len(lstack) > 0 {
				// Add as child to previous item's content
				parentItem := lstack[len(lstack)-1].item
				if parentItem != nil {
					parentItem.Content = append(parentItem.Content, ln)
				}
			} else {
				rootList = ln
			}
			item = &ListItemNode{}
			ln.Items = append(ln.Items, *item)
			lstack = append(lstack, &listStackElem{level: level, ordered: ordered, list: ln, item: item})
		} else {
			// Same level/type: start new item
			ln := lstack[len(lstack)-1].list
			item := &ListItemNode{}
			ln.Items = append(ln.Items, *item)
			lstack[len(lstack)-1].item = item
		}
		// Add content as TextNode to current item
		curItem := lstack[len(lstack)-1].item
		if curItem != nil && strings.TrimSpace(content) != "" {
			curItem.Content = append(curItem.Content, &TextNode{Text: strings.TrimSpace(content)})
		}
	}
	// Unwind stack and add rootList to parent
	if rootList == nil && len(lstack) > 0 {
		rootList = lstack[0].list
	}
	if rootList != nil {
		if node, ok := parent.(interface{ AddChild(BlockNode) }); ok {
			node.AddChild(rootList)
		}
		*stack = append(*stack, rootList)
		return true
	}
	return false
}
