package xatena

import (
	"testing"
)

const stoppTestData = `
=== test
--- input
><blockquote>
<p>test</p>
</blockquote><
--- expected
<blockquote>
<p>test</p>
</blockquote>

=== test
--- input
><ins datetime="2010-02-27T00:00:00Z"><
foobar
></ins><
--- expected
<ins datetime="2010-02-27T00:00:00Z">
<p>foobar</p>
</ins>


`

func TestFormat_StopP_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(stoppTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
