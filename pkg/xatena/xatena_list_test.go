package xatena

import (
	"testing"
)

const listTestData = `=== unordered_list
--- input
- ul
- ul
-- ul
-- ul
--- ul
- ul
--- expected
<ul>
  <li>ul</li>
  <li>ul</li>
  <li>
    <ul>
      <li>ul</li>
      <li>ul</li>
      <li>
        <ul>
          <li>ul</li>
        </ul>
      </li>
    </ul>
  </li>
  <li>ul</li>
</ul>

=== ordered_list
--- input
+ ol
+ ol
++ ol
++ ol
+++ ol
+ ol
--- expected
<ol>
  <li>ol</li>
  <li>ol</li>
  <li>
    <ol>
      <li>ol</li>
      <li>ol</li>
      <li>
        <ol>
          <li>ol</li>
        </ol>
      </li>
    </ol>
  </li>
  <li>ol</li>
</ol>
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
