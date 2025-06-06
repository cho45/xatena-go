package syntax

import (
	"context"
	htmltpl "html/template"
	"regexp"
	"strings"
)

var TableTemplate = htmltpl.Must(htmltpl.New("table").Parse(`
<table>
{{- range .Rows}}
  <tr>
  {{- range .}}
    {{if .IsHeader}}<th>{{.Content}}</th>{{else}}<td>{{.Content}}</td>{{end}}
  {{- end}}
  </tr>
{{- end}}
</table>`))

// TableNode represents a table block
type TableNode struct {
	Rows [][]TableCellNode
}

type TableCellNode struct {
	IsHeader bool
	Content  string
}

func (t *TableNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	type cell struct {
		IsHeader bool
		Content  htmltpl.HTML
	}
	type row []cell
	var rows []row
	inline := xatena.GetInline()
	for _, r := range t.Rows {
		var rowCells row
		for _, c := range r {
			rowCells = append(rowCells, cell{
				IsHeader: c.IsHeader,
				Content:  htmltpl.HTML(inline.Format(ctx, c.Content)),
			})
		}
		rows = append(rows, rowCells)
	}
	params := map[string]interface{}{
		"Rows": rows,
	}
	html := xatena.ExecuteTemplate("table", params)
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

func (p *TableParser) CanHandle(line string) bool {
	return strings.HasPrefix(line, "|")
}

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
	parent.AddChild(node)
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
