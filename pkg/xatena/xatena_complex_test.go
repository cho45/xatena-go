package xatena

import (
	"testing"
)

const complexTestData = `
=== test
--- input
* This is a head

foobar

barbaz

:foo:bar
:foo:bar

- list1
- list1
- list1

>|perl|
test code
||<

ok?
--- expected
<div class="section">
    <h3>This is a head</h3>
    <p>foobar</p>
    <p>barbaz</p>
    <dl>
        <dt>foo</dt>
        <dd>bar</dd>
        <dt>foo</dt>
        <dd>bar</dd>
    </dl>
    <ul>
        <li>list1</li>
        <li>list1</li>
        <li>list1</li>
    </ul>
    <pre class="code lang-perl">test code</pre>
    <p>ok?</p>
</div>


=== test
--- input
>||
<!--
test
-->
||<

--- expected
<pre class="code">
&lt;!--
test
--&gt;
</pre>

=== test
--- input
<!--

>||
test
||<

-->

--- expected
<!-- -->


`

func TestFormat_ComplexIntegration(t *testing.T) {
	blocks := parseTestBlocks(complexTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
