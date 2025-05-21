package syntax

import (
	"context"
	"regexp"
	"strings"

	"github.com/cho45/xatena-go/internal/util"
)

// TableNode represents a table block
type TableNode struct {
	Rows [][]TableCellNode
}

type TableCellNode struct {
	IsHeader bool
	Content  string
}

func (t *TableNode) ToHTML(ctx context.Context, inline Inline) string {
	html := "<table>\n"
	for _, row := range t.Rows {
		html += "  <tr>\n"
		for _, cell := range row {
			if cell.IsHeader {
				html += "    <th>" + util.EscapeHTML(cell.Content) + "</th>\n"
			} else {
				html += "    <td>" + util.EscapeHTML(cell.Content) + "</td>\n"
			}
		}
		html += "  </tr>\n"
	}
	html += "</table>"
	return html
}

func (t *TableNode) AddChild(n Node) {
	// テーブルは子ブロックを持たないので何もしない
}

func (t *TableNode) GetContent() []Node {
	return nil
}

type TableParser struct{}

var reTableRow = regexp.MustCompile(`^\|`)

func (p *TableParser) Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool {
	if !reTableRow.MatchString(scanner.Peek()) {
		return false
	}
	var rows [][]TableCellNode
	for !scanner.EOF() && reTableRow.MatchString(scanner.Peek()) {
		rows = append(rows, parseTableRow(scanner.Next()))
	}
	if len(rows) == 0 {
		return false
	}

	node := &TableNode{Rows: rows}
	if add, ok := parent.(interface{ AddChild(Node) }); ok {
		add.AddChild(node)
	}
	*stack = append(*stack, node)
	return true
}

func parseTableRow(line string) []TableCellNode {
	// |a|b|c| → [a b c]
	var cells []TableCellNode
	// 先頭と末尾の | を除去
	trimmed := strings.Trim(line, "| ")
	parts := strings.Split(trimmed, "|")
	for _, cell := range parts {
		cell = strings.TrimSpace(cell)
		isHeader := len(cell) > 0 && cell[0] == '*'
		content := cell
		if isHeader {
			content = cell[1:]
		}
		cells = append(cells, TableCellNode{IsHeader: isHeader, Content: content})
	}
	return cells
}
