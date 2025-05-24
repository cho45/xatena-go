package xatena

import (
	"testing"
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
