package xatena

import (
	"testing"
)

const preTestData = `
=== test
--- input
>|
quote
|<
--- expected
<pre>
quote
</pre>

=== test
--- input
>|
quote1
>>
quote2
<<
|<
--- expected
<pre>
quote1
    <blockquote>
        quote2
    </blockquote>
</pre>

=== test
--- input
>|
http://example.com/
|<
--- expected
<pre>
<a href="http://example.com/">http://example.com/</a>
</pre>

=== test
--- input
>|
quote
|<
test
--- expected
<pre>
quote
</pre>
<p>test</p>


`

func TestFormat_Pre_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(preTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
