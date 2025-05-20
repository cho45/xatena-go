package xatena

import (
	"testing"
)

const tableTestData = `=== table_simple
--- input
|*foo|*bar|*baz|
|test|test|test|
|test|test|test|
--- expected
<table>
  <tr>
    <th>foo</th>
    <th>bar</th>
    <th>baz</th>
  </tr>
  <tr>
    <td>test</td>
    <td>test</td>
    <td>test</td>
  </tr>
  <tr>
    <td>test</td>
    <td>test</td>
    <td>test</td>
  </tr>
</table>
`

func TestFormat_Table_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(tableTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
