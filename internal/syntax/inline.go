package syntax

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type Inline interface {
	Format(s string) string
}

type InlineFormatter struct {
	footnotes  []Footnote
	aggressive bool              // aggressive title fetch (未実装)
	cache      map[string]string // dummy cache for title (未実装)
}

type Footnote struct {
	Number int
	Note   string
	Title  string
}

func NewInlineFormatter() *InlineFormatter {
	return &InlineFormatter{
		footnotes: []Footnote{},
		cache:     map[string]string{},
	}
}

func (f *InlineFormatter) Format(s string) string {
	s = strings.TrimPrefix(s, "\n")

	// 個別の記法ごとの正規表現とハンドラ
	type inlineRule struct {
		pattern *regexp.Regexp
		handler func([]string) string
	}

	rules := []inlineRule{
		{
			regexp.MustCompile(`\[\]([\s\S]*?)\[\]`),
			func(m []string) string { return m[1] }, // [[]...[]] アンリンク
		},
		{
			regexp.MustCompile(`\(\(\(.*?\)\)\)`),
			func(m []string) string { return m[0][1 : len(m[0])-1] }, // (((...)))
		},
		{
			regexp.MustCompile(`\)\(\(.*?\)\)\(`),
			func(m []string) string { return m[0][1 : len(m[0])-1] }, // )((...))(
		},
		{
			regexp.MustCompile(`\(\((.+?)\)\)`),
			func(m []string) string { // footnote
				note := m[1]
				title := stripTags(note)
				f.footnotes = append(f.footnotes, Footnote{Number: len(f.footnotes) + 1, Note: note, Title: title})
				return fmt.Sprintf(`<a href="#fn%d" title="%s">*%d</a>`, len(f.footnotes), html.EscapeString(title), len(f.footnotes))
			},
		},
		{
			regexp.MustCompile(`(?i)<a[^>]+>[\s\S]*?</a>`),
			func(m []string) string { return m[0] }, // <a...>...</a>
		},
		{
			regexp.MustCompile(`<!--.*?-->`),
			func(m []string) string { return "<!-- -->" }, // コメント
		},
		{
			regexp.MustCompile(`(?i)<[^>]+>`),
			func(m []string) string { return m[0] }, // その他タグ
		},
		{
			regexp.MustCompile(`\[((?:https?|ftp)://[^\s:]+(?:\:\d+)?[^\s:]+)(:(?:title(?:=([^\]]+))?|barcode))?\]`),
			func(m []string) string { // [url:option]
				uri, opt, title := m[1], m[2], m[3]
				if opt == ":barcode" {
					return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?chs=150x150&cht=qr&chl=%s" title="%s"/>`, url.QueryEscape(uri), html.EscapeString(uri))
				}
				if strings.HasPrefix(opt, ":title") {
					return fmt.Sprintf(`<a href="%s">%s</a>`, uri, html.EscapeString(title))
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
		{
			regexp.MustCompile(`\[((?:https?|ftp):[^\s<>\]]+)\]`),
			func(m []string) string { // [url]
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
		{
			regexp.MustCompile(`\[mailto:([^\s\@:?]+\@[^
\s\@:?]+(\?[^\s]+)?)\]`),
			func(m []string) string { // mailto
				uri := m[1]
				return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, uri, uri)
			},
		},
		{
			regexp.MustCompile(`\[tex:([^\]]+)\]`),
			func(m []string) string { // tex
				tex := m[1]
				return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?cht=tx&chl=%s" alt="%s"/>`, url.QueryEscape(tex), html.EscapeString(tex))
			},
		},
		// 裸URL検出用パターンを追加
		{
			regexp.MustCompile(`((?:https?|ftp):[^\s<>"]+)`),
			func(m []string) string {
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
	}

	// すべてのパターンを | で結合した大きな正規表現を作成
	var patterns []string
	for _, r := range rules {
		patterns = append(patterns, r.pattern.String())
	}
	bigRe := regexp.MustCompile(strings.Join(patterns, "|"))

	// 1パスで置換
	result := bigRe.ReplaceAllStringFunc(s, func(m string) string {
		for _, r := range rules {
			if sub := r.pattern.FindStringSubmatch(m); sub != nil {
				return r.handler(sub)
			}
		}
		return m
	})
	return result
}

func stripTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}

func (f *InlineFormatter) Footnotes() []Footnote {
	return f.footnotes
}
