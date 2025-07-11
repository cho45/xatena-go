package xatena

import (
	"context"
	"testing"

	"github.com/cho45/xatena-go/internal/syntax"
)

const tableTestData = `
=== test
--- input
|*head|*head|*head|
|foo|bar|baz|
|foo|bar|baz|
--- expected
<table>
    <tr>
        <th>head</th>
        <th>head</th>
        <th>head</th>
    </tr>
    <tr>
        <td>foo</td>
        <td>bar</td>
        <td>baz</td>
    </tr>
    <tr>
        <td>foo</td>
        <td>bar</td>
        <td>baz</td>
    </tr>
</table>

=== test
--- input
|*head|*head|*head|
|http://www.lowreal.net/|bar|baz|
--- expected
<table>
    <tr>
        <th>head</th>
        <th>head</th>
        <th>head</th>
    </tr>
    <tr>
        <td><a href="http://www.lowreal.net/">http://www.lowreal.net/</a></td>
        <td>bar</td>
        <td>baz</td>
    </tr>
</table>

=== test
--- input
|*head|*head|*head|
|foo|bar|baz|
|foo|bar|baz|
test
--- expected
<table>
    <tr>
        <th>head</th>
        <th>head</th>
        <th>head</th>
    </tr>
    <tr>
        <td>foo</td>
        <td>bar</td>
        <td>baz</td>
    </tr>
    <tr>
        <td>foo</td>
        <td>bar</td>
        <td>baz</td>
    </tr>
</table>
<p>test</p>


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

// TestTableNodeHasContent tests AddChild and GetContent methods for TableNode
func TestTableNodeHasContent(t *testing.T) {
	x := NewXatena()
	
	// Parse a table to get TableNode
	input := "|cell1|cell2|\n|data1|data2|"
	root := x.parseXatena(context.Background(), input)
	
	if len(root.Content) == 0 {
		t.Fatal("expected at least one node in parsed content")
	}
	
	// Find the TableNode
	var tableNode *syntax.TableNode
	for _, node := range root.Content {
		if tn, ok := node.(*syntax.TableNode); ok {
			tableNode = tn
			break
		}
	}
	
	if tableNode == nil {
		t.Fatal("expected to find a TableNode")
	}
	
	// Test AddChild method (should not panic, but does nothing)
	textNode := &syntax.TextNode{Text: "test"}
	tableNode.AddChild(textNode)
	
	// Test GetContent method (should not panic, returns nil)
	content := tableNode.GetContent()
	if content != nil {
		t.Errorf("expected nil content from TableNode.GetContent(), got %v", content)
	}
}
