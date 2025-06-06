package xatena

import (
	"testing"
)

const seeMoreTestData = `
### test
::: input
foobar

>||
superpre
====
||<

====

barbaz

** head

*** head

foo

** head

bar

::: expected
<p>foobar</p>
<pre class="code">superpre
====</pre>
<div class="seemore">
	<p>barbaz</p>
	<div class="section">
		<h4>head</h4>
		<div class="section">
			<h5>head</h5>
			<p>foo</p>
		</div>
	</div>
	<div class="section">
		<h4>head</h4>
		<p>bar</p>
	</div>
</div>

### test
::: input
* head

foobar

====

barbaz

* head

foo

::: expected
<div class="section">
	<h3>head</h3>
	<p>foobar</p>

	<div class="seemore">
		<p>barbaz</p>
	</div>
</div>

<div class="section">
	<h3>head</h3>
	<p>foo</p>
</div>

### super seemore
::: input
* head

foobar

=====

barbaz

* head

foo

::: expected
<div class="section">
	<h3>head</h3>
	<p>foobar</p>

	<div class="seemore">
		<p>barbaz</p>

		<div class="section">
			<h3>head</h3>
			<p>foo</p>
		</div>
	</div>
</div>

### super seemore
::: input
* head

foobar

=====
* head

foo

::: expected
<div class="section">
	<h3>head</h3>
	<p>foobar</p>

	<div class="seemore">
		<div class="section">
			<h3>head</h3>
			<p>foo</p>
		</div>
	</div>
</div>


`

func TestFormat_SeeMore_ENDStyle(t *testing.T) {
	blocks := parseTestBlocksWithDelim(seeMoreTestData, "###", ":::")
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
