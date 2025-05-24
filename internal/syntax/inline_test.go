package syntax

import (
	"testing"
)

func TestInlineFormatter_Format(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "simple link",
			input:  "[http://example.com/]",
			expect: `<a href="http://example.com/">http://example.com/</a>`,
		},
		{
			name:   "simple link raw",
			input:  "http://example.com/",
			expect: `<a href="http://example.com/">http://example.com/</a>`,
		},
		{
			name:   "mailto",
			input:  "[mailto:foo@example.com]",
			expect: `<a href="mailto:foo@example.com">foo@example.com</a>`,
		},
		{
			name:   "tex",
			input:  "[tex:E=mc^2]",
			expect: `<img src="http://chart.apis.google.com/chart?cht=tx&chl=E%3Dmc%5E2" alt="E=mc^2"/>`,
		},
		{
			name:   "footnote",
			input:  "((note))",
			expect: `<a href="#fn1" title="note">*1</a>`,
		},
		{
			name:   "barcode",
			input:  "[http://example.com/:barcode]",
			expect: `<img src="http://chart.apis.google.com/chart?chs=150x150&cht=qr&chl=http%3A%2F%2Fexample.com%2F" title="http://example.com/"/>`,
		},
		{
			name:   "title option",
			input:  "[http://example.com/:title=Example]",
			expect: `<a href="http://example.com/">Example</a>`,
		},
		{
			name:   "unlink",
			input:  "[]unlinked[]",
			expect: "unlinked",
		},
	}

	for _, c := range cases {
		f := NewInlineFormatter()
		got := f.Format(c.input)
		if got != c.expect {
			t.Errorf("%s: input=%q\nexpect=%q\ngot   =%q", c.name, c.input, c.expect, got)
		}
	}
}
