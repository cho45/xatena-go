package xatena

import (
	"testing"
)

const paragraphTestData = `
=== test
--- input
test
--- expected
<p>test</p>

=== test
--- input
test
test
--- expected
<p>test<br />test</p>

=== test
--- input
test
test

test
--- expected
<p>test<br />test</p>
<p>test</p>

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

func TestFormat_Paragraph_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(paragraphTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
