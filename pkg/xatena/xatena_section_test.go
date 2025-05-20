package xatena

import (
	"testing"
)

const sectionTestData = `=== section_simple
--- input
* head1
foobar
** head2
*** head3
--- expected
<div class="section">
<h3>head1</h3>
<p>foobar</p>
  <div class="section">
  <h4>head2</h4>
    <div class="section">
    <h5>head3</h5>
    </div>
  </div>
</div>
`

func TestFormat_Section_ENDStyle(t *testing.T) {
	blocks := parseTestBlocks(sectionTestData)
	for _, b := range blocks {
		input := b.Sections["input"]
		expected := b.Sections["expected"]
		t.Run(b.Name, func(t *testing.T) {
			got := Format(input)
			EqualHTML(t, got, expected)
		})
	}
}
