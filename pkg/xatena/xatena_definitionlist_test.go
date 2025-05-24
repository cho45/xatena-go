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

=== test
--- input
:foo:bar
:baz:piyo
--- expected
<dl>
    <dt>foo</dt>
    <dd>bar</dd>
    <dt>baz</dt>
    <dd>piyo</dd>
</dl>

=== test
--- input
:foo:http://www.lowreal.net/
:baz:piyo
--- expected
<dl>
    <dt>foo</dt>
    <dd><a href="http://www.lowreal.net/">http://www.lowreal.net/</a></dd>
    <dt>baz</dt>
    <dd>piyo</dd>
</dl>

=== test
--- input
:foo:http://www.lowreal.net/
:baz:piyo
--- expected
<dl>
    <dt>foo</dt>
    <dd><a href="http://www.lowreal.net/">http://www.lowreal.net/</a></dd>
    <dt>baz</dt>
    <dd>piyo</dd>
</dl>

=== test
--- input
:foo:
::http://www.lowreal.net/
:baz:
::piyo
::piyo
--- expected
<dl>
    <dt>foo</dt>
    <dd><a href="http://www.lowreal.net/">http://www.lowreal.net/</a></dd>
    <dt>baz</dt>
    <dd>piyo</dd>
    <dd>piyo</dd>
</dl>

=== test
--- input
:foo
--- expected
<p>:foo</p>

=== test
--- input
:foo:bar
:baz:piyo
test
--- expected
<dl>
    <dt>foo</dt>
    <dd>bar</dd>
    <dt>baz</dt>
    <dd>piyo</dd>
</dl>
<p>test</p>

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
