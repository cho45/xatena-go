package xatena

import (
	"context"
	"testing"

	"github.com/cho45/xatena-go/internal/syntax"
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

// TestDefinitionListNodeHasContent tests AddChild and GetContent methods for DefinitionListNode
func TestDefinitionListNodeHasContent(t *testing.T) {
	x := NewXatena()
	
	// Parse a definition list to get DefinitionListNode
	input := ":term:description"
	root := x.parseXatena(context.Background(), input)
	
	if len(root.Content) == 0 {
		t.Fatal("expected at least one node in parsed content")
	}
	
	// Find the DefinitionListNode
	var definitionListNode *syntax.DefinitionListNode
	for _, node := range root.Content {
		if dln, ok := node.(*syntax.DefinitionListNode); ok {
			definitionListNode = dln
			break
		}
	}
	
	if definitionListNode == nil {
		t.Fatal("expected to find a DefinitionListNode")
	}
	
	// Test AddChild method (should not panic, but does nothing)
	textNode := &syntax.TextNode{Text: "test"}
	definitionListNode.AddChild(textNode)
	
	// Test GetContent method (should not panic, returns nil)
	content := definitionListNode.GetContent()
	if content != nil {
		t.Errorf("expected nil content from DefinitionListNode.GetContent(), got %v", content)
	}
}
