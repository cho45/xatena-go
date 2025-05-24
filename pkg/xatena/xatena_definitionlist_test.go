package xatena

import (
	"testing"
)

const definitionListTestData = `=== definitionlist_simple
--- input
:foo:bar
:baz:qux
--- expected
<dl>
  <dt>foo</dt>
  <dd>bar</dd>
  <dt>baz</dt>
  <dd>qux</dd>
</dl>

=== definitionlist_multiline
--- input
:foo:
:: bar
:baz:qux
--- expected
<dl>
  <dt>foo</dt>
  <dd> bar</dd>
  <dt>baz</dt>
  <dd>qux</dd>
</dl>
`

func TestFormat_DefinitionList_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(definitionListTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
