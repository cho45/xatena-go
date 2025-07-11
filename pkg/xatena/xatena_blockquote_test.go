package xatena

import (
	"testing"
)

const blockquoteTestData = `
=== test
--- input
>>
quote
<<
--- expected
<blockquote>
<p>quote</p>
</blockquote>

=== test
--- input
>>
quote1
>>
quote2
<<
<<
--- expected
<blockquote>
    <p>quote1</p>
    <blockquote>
        <p>quote2</p>
    </blockquote>
</blockquote>

=== test
--- input
>http://example.com/>
quote
<<
--- expected
<blockquote cite="http://example.com/">
	<p>quote</p>
	<cite><a href="http://example.com/">http://example.com/</a></cite>
</blockquote>

=== http
--- input
>http://example.com/:title>
quote
<<
--- expected
<blockquote cite="http://example.com/">
	<p>quote</p>
	<cite><a href="http://example.com/">Example Web Page</a></cite>
</blockquote>

=== http
--- input
>http://example.com/:title=TEST>
quote
<<
--- expected
<blockquote cite="http://example.com/">
	<p>quote</p>
	<cite><a href="http://example.com/">TEST</a></cite>
</blockquote>

=== cite
--- input
>foobar>
quote
<<
--- expected
<blockquote>
	<p>quote</p>
	<cite>foobar</cite>
</blockquote>


=== test
--- input
>>
quote
<<
test
--- expected
<blockquote>
<p>quote</p>
</blockquote>
<p>test</p>

=== bug1
--- input
>>
* hoge1
hoge2
<<
hoge3
--- expected
<blockquote>
	<div class="section">
		<h3>hoge1</h3>
		<p>hoge2</p>
	</div>
</blockquote>
<p>hoge3</p>

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
