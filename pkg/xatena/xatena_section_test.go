package xatena

import (
	"testing"

	"github.com/cho45/xatena-go/internal/syntax"
)

const sectionTestData = `
=== test
--- input
* test1

foo

* test2

bar
--- expected
<div class="section">
<h3>test1</h3>
<p>foo</p>
</div>

<div class="section">
<h3>test2</h3>
<p>bar</p>
</div>


=== test
--- input
* test1

foo

** test1.1

foo!

** test1.2

foo!

*** test1.2.1

foo!

* test2

bar
--- expected
<div class="section">
<h3>test1</h3>
<p>foo</p>

<div class="section">
<h4>test1.1</h4>
<p>foo!</p>
</div>

<div class="section">
<h4>test1.2</h4>
<p>foo!</p>

<div class="section">
<h5>test1.2.1</h5>
<p>foo!</p>
</div>
</div>
</div>

<div class="section">
<h3>test2</h3>
<p>bar</p>
</div>


=== test
--- input
* http://example.com/
foo
--- expected
<div class="section">
<h3><a href="http://example.com/">http://example.com/</a></h3>
<p>foo</p>
</div>

=== heading
--- input
* ***
foobar
--- expected
<div class="section">
	<h3>***</h3>
	<p>foobar</p>
</div>

=== very complex heading
--- input
****foobar
foobar
--- expected
<div class="section">
	<h3>***foobar</h3>
	<p>foobar</p>
</div>

=== very complex heading
--- input
* 
trailing space

** 
trailing space

*** 
trailing space

*
no spaces

**
no spaces

***
no spaces
--- expected
<div class="section">
	<h3></h3>
	<p>trailing space</p>
	<div class="section">
		<h4></h4>
		<p>trailing space</p>

		<div class="section">
			<h5></h5>
			<p>trailing space</p>
		</div>
	</div>
</div>

<div class="section">
	<h3></h3>
	<p>no spaces</p>
</div>

<div class="section">
	<h3>*</h3>
	<p>no spaces</p>
</div>

<div class="section">
	<h3>**</h3>
	<p>no spaces</p>
</div>


`

func TestFormat_Section_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(sectionTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}

// TestSectionTitleNodeGetContent tests GetContent method for SectionTitleNode
func TestSectionTitleNodeGetContent(t *testing.T) {
	// Create a SectionTitleNode directly
	titleNode := &syntax.SectionTitleNode{}
	
	// Test GetContent method (should not panic, returns nil)
	content := titleNode.GetContent()
	if content != nil {
		t.Errorf("expected nil content from SectionTitleNode.GetContent(), got %v", content)
	}
}
