package xatena

import (
	"testing"
)

const listTestData = `
=== test
--- input
- 1
- 2
- 3
--- expected
<ul>
    <li>1</li>
    <li>2</li>
    <li>3</li>
</ul>

=== test
--- input
- 1
- 2
-- 2.1
-- 2.2
--+ 2.2.3
- 3
--- expected
<ul>
    <li>1</li>
    <li>2
        <ul>
            <li>2.1</li>
            <li>2.2
                <ol>
                    <li>2.2.3</li>
                </ol>
            </li>
        </ul>
    </li>
    <li>3</li>
</ul>

=== test
--- input
- http://www.lowreal.net/
- 2
-+ 2.1
-+ 2.2
- 3
--- expected
<ul>
    <li><a href="http://www.lowreal.net/">http://www.lowreal.net/</a></li>
    <li>2
        <ol>
            <li>2.1</li>
            <li>2.2</li>
        </ol>
    </li>
    <li>3</li>
</ul>

=== test
--- input
 -foo
--- expected
<p>-foo</p>

=== test
--- input
- 1
- 2
- 3
test
--- expected
<ul>
    <li>1</li>
    <li>2</li>
    <li>3</li>
</ul>
<p>test</p>

=== bug
--- input
foo

-

bar
--- expected
<p>foo</p>
<p>-</p>
<p>bar</p>

=== bug
--- input
++ unko

--- expected
<ol>
	<li>
		<ol>
			<li>unko</li>
		</ol>
	</li>
</ol>

=== bug
--- input
-- unko

--- expected
<ul>
	<li>
		<ul>
			<li>unko</li>
		</ul>
	</li>
</ul>

=== bug
--- input
+++ unko

--- expected
<ol>
	<li>
		<ol>
			<li>
				<ol>
					<li>unko</li>
				</ol>
			</li>
		</ol>
	</li>
</ol>

=== bug
--- input
+ foo
- bar

--- expected
<ol>
	<li>foo</li>
</ol>
<ul>
	<li>bar</li>
</ul>

`

func TestFormat_List_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(listTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
