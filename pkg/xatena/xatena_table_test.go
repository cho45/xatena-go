package xatena

import (
	"testing"
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
