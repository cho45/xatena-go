package xatena

import (
	"testing"
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
