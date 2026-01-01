package xatena

import (
	"context"
	"strings"
	"testing"
)

func TestBugs_ReportedCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []string // 期待される部分文字列
		reject []string // 含まれてはいけない部分文字列
	}{
		{
			name:  "Bug ID 13392: Unclosed comment should still be treated as comment and closed in HTML",
			input: "text before\n><!--\ncomment body",
			// 1. "text before" は出力されるべき
			// 2. 閉じのない <!-- 以降は（現状の仕様では）コメントとして処理される
			// 3. ただし、生成される HTML はちゃんと閉じられているべき
			want: []string{"text before", "<!--", "-->"},
		},
		{
			name:  "Bug ID 13468: HTML tags in footnotes should be stripped in title attribute",
			input: "((footnote with <em>tags</em>))",
			// title="..." should have tags stripped
			want:   []string{`title="footnote with tags"`},
			reject: []string{"<em>", "&lt;em&gt;"},
		},
		{
			name:  "Bug ID 13275: Isolated closing notation should not cause subsequent content to be lost",
			input: "part A\n|<\npart B",
			// パニック回避後、孤立した |< があっても前後のテキストが失われないこと
			want: []string{"part A", "part B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewXatena()
			got := x.ToHTML(context.Background(), tt.input)
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("%s: output does not contain %q\nGot: %s", tt.name, w, got)
				}
			}
			for _, r := range tt.reject {
				if strings.Contains(got, r) {
					t.Errorf("%s: output contains rejected string %q\nGot: %s", tt.name, r, got)
				}
			}
		})
	}
}
