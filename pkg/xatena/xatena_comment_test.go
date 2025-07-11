package xatena

import (
	"context"
	"testing"

	"github.com/cho45/xatena-go/internal/syntax"
)

const commentTestData = `
=== test
--- input
<!--
secret
-->

--- expected
<!-- -->

=== test
--- input
<!-- secret -->

--- expected
<!-- -->

=== test
--- input
<!-- secret -->
foobar

--- expected
<!-- -->
<p>foobar</p>

=== test
--- input
foobar <!-- secret -->

--- expected
<p>foobar</p>
<!-- -->

=== test
--- input
- <!-- foobar -->
- 1

--- expected
<ul>
    <li><!-- --></li>
    <li>1</li>
</ul>

=== test inline comment
--- input
- baz <!-- foobar --> bar
- 1

--- expected
<ul>
    <li>baz <!-- --> bar</li>
    <li>1</li>
</ul>

=== test inline comment
--- input
>|
test
<!-- foobar -->
test
|<

--- expected
<pre>
test
<!-- -->
test
</pre>


`

func TestFormat_Comment(t *testing.T) {
	blocks := parseTestBlocks(commentTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}

// TestCommentNodeHasContent tests AddChild and GetContent methods for CommentNode
func TestCommentNodeHasContent(t *testing.T) {
	// CommentNode doesn't actually use AddChild/GetContent in practice,
	// but we test the interface methods for coverage
	x := NewXatena()
	
	// Parse a comment to get CommentNode
	input := "<!--\ntest content\n-->"
	root := x.parseXatena(context.Background(), input)
	
	if len(root.Content) == 0 {
		t.Fatal("expected at least one node in parsed content")
	}
	
	// Find the CommentNode
	var commentNode *syntax.CommentNode
	for _, node := range root.Content {
		if cn, ok := node.(*syntax.CommentNode); ok {
			commentNode = cn
			break
		}
	}
	
	if commentNode == nil {
		t.Fatal("expected to find a CommentNode")
	}
	
	// Test AddChild method (should not panic, but does nothing)
	textNode := &syntax.TextNode{Text: "test"}
	commentNode.AddChild(textNode)
	
	// Test GetContent method (should not panic, returns nil)
	content := commentNode.GetContent()
	if content != nil {
		t.Errorf("expected nil content from CommentNode.GetContent(), got %v", content)
	}
}
