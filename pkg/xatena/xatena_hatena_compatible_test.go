package xatena

import (
	"context"
	"testing"
)

const hatenaCompatibleTestData = `

=== test
--- input
test
--- expected
<p>test</p>

=== test
--- input
test http://example.com/
test
--- expected
<p>test <a href="http://example.com/">http://example.com/</a></p>
<p>test</p>

=== test
--- input
test
test

test
--- expected
<p>test</p>
<p>test</p>
<p>test</p>

=== test
--- input
* test1

foo

* test2

bar
--- expected
<h3>test1</h3>
<p>foo</p>

<h3>test2</h3>
<p>bar</p>


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
<h3>test1</h3>
<p>foo</p>

<h4>test1.1</h4>
<p>foo!</p>

<h4>test1.2</h4>
<p>foo!</p>

<h5>test1.2.1</h5>
<p>foo!</p>

<h3>test2</h3>
<p>bar</p>


=== test
--- input
* http://example.com/
foo
--- expected
<h3><a href="http://example.com/">http://example.com/</a></h3>
<p>foo</p>

=== test
--- input
a


a
--- expected
<p>a</p>
<br />
<p>a</p>

=== test
--- input
a



a
--- expected
<p>a</p>
<br />
<br />
<p>a</p>

=== test
--- input
a




a
--- expected
<p>a</p>
<br />
<br />
<br />
<p>a</p>


`

func TestFormat_HatenaCompatible_ENDStyle(t *testing.T) {
	x := NewXatenaWithFields(NewInlineFormatter(), true)

	blocks := parseTestBlocks(hatenaCompatibleTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := x.ToHTML(context.Background(), input)
			EqualHTML(t, got, expected)
		})
	}
}
