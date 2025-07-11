package xatena

import (
	"context"
	"testing"

	"github.com/cho45/xatena-go/internal/syntax"
)

const superPreTestData = `
=== test
--- input
>||
test
test
||<
--- expected
<pre class="code">
test
test
</pre>

=== test
--- input
>||
<a href="foobar">foobar</a>
||<
--- expected
<pre class="code">
&lt;a href=&#34;foobar&#34;&gt;foobar&lt;/a&gt;
</pre>

=== test
--- input
>||
>>
foobar
<<
||<
--- expected
<pre class="code">
&gt;&gt;
foobar
&lt;&lt;
</pre>

=== test
--- input
>|perl|
test
test
||<
--- expected
<pre class="code lang-perl">
test
test
</pre>

=== test
--- input
>|perl|
test
test
||,
--- expected
<pre class="code lang-perl">
test
test
</pre>

=== test
--- input
>||
test
test
||<
test
--- expected
<pre class="code">
test
test
</pre>
<p>test</p>


`

func TestFormat_SuperPre_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(superPreTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}

// TestSuperPreNodeHasContent tests AddChild and GetContent methods for SuperPreNode
func TestSuperPreNodeHasContent(t *testing.T) {
	x := NewXatena()
	
	// Parse a superpre to get SuperPreNode
	input := ">||\ncode\n||<"
	root := x.parseXatena(context.Background(), input)
	
	if len(root.Content) == 0 {
		t.Fatal("expected at least one node in parsed content")
	}
	
	// Find the SuperPreNode
	var superPreNode *syntax.SuperPreNode
	for _, node := range root.Content {
		if spn, ok := node.(*syntax.SuperPreNode); ok {
			superPreNode = spn
			break
		}
	}
	
	if superPreNode == nil {
		t.Fatal("expected to find a SuperPreNode")
	}
	
	// Test AddChild method (should not panic, but does nothing)
	textNode := &syntax.TextNode{Text: "test"}
	superPreNode.AddChild(textNode)
	
	// Test GetContent method (should not panic, returns nil)
	content := superPreNode.GetContent()
	if content != nil {
		t.Errorf("expected nil content from SuperPreNode.GetContent(), got %v", content)
	}
}
