package xatena

import (
	"testing"
)

const blockquoteTestData = `=== blockquote_simple
--- input
>>
quoted text
foobar
<<
--- expected
<blockquote>
<p>quoted text</p>
<p>foobar</p>
</blockquote>

=== blockquote_with_cite
--- input
>http://example.com/>
foobar
<<
--- expected
<blockquote cite="http://example.com/">
<p>foobar</p>
<cite><a href="http://example.com/">http://example.com/</a></cite>
</blockquote>
`

func TestFormat_Blockquote_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(blockquoteTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
