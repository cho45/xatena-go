package xatena

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type InlineRule struct {
	Pattern *regexp.Regexp
	Handler func(f *InlineFormatter, m []string) string
}

type InlineFormatter struct {
	footnotes []Footnote
	rules     []InlineRule
}

type Footnote struct {
	Number int
	Note   string
	Title  string
}

func (f *InlineFormatter) AddRule(rule InlineRule) {
	f.rules = append(f.rules, rule)
}

func (f *InlineFormatter) AddRuleAt(index int, rule InlineRule) {
	if index < 0 || index > len(f.rules) {
		index = len(f.rules)
	}
	f.rules = append(f.rules[:index], append([]InlineRule{rule}, f.rules[index:]...)...)
}

func NewInlineFormatter() *InlineFormatter {
	f := &InlineFormatter{
		footnotes: []Footnote{},
	}
	// デフォルトルールを登録
	f.rules = defaultInlineRules()
	return f
}

// デフォルトルール群を返す
func defaultInlineRules() []InlineRule {
	return []InlineRule{
		{
			Pattern: regexp.MustCompile(`\[\]([\s\S]*?)\[\]`),
			Handler: func(f *InlineFormatter, m []string) string { return m[1] }, // [[]...[]] アンリンク
		},
		{
			Pattern: regexp.MustCompile(`\(\(\(.*?\)\)\)`),
			Handler: func(f *InlineFormatter, m []string) string { return m[0][1 : len(m[0])-1] }, // (((...)))
		},
		{
			Pattern: regexp.MustCompile(`\)\(\(.*?\)\)\(`),
			Handler: func(f *InlineFormatter, m []string) string { return m[0][1 : len(m[0])-1] }, // )((...))(
		},
		{
			Pattern: regexp.MustCompile(`\(\((.+?)\)\)`),
			Handler: func(f *InlineFormatter, m []string) string { // footnote
				note := m[1]
				title := html.EscapeString(note)
				f.footnotes = append(f.footnotes, Footnote{Number: len(f.footnotes) + 1, Note: note, Title: title})
				return fmt.Sprintf(`<a href="#fn%d" title="%s">*%d</a>`, len(f.footnotes), html.EscapeString(title), len(f.footnotes))
			},
		},
		{
			Pattern: regexp.MustCompile(`(?i)<a[^>]+>[\s\S]*?</a>`),
			Handler: func(f *InlineFormatter, m []string) string { return m[0] }, // <a...>...</a>
		},
		{
			Pattern: regexp.MustCompile(`<!--.*?-->`),
			Handler: func(f *InlineFormatter, m []string) string { return "<!-- -->" }, // コメント
		},
		{
			Pattern: regexp.MustCompile(`(?i)<[^>]+>`),
			Handler: func(f *InlineFormatter, m []string) string { return m[0] }, // その他タグ
		},
		{
			Pattern: regexp.MustCompile(`\[((?:https?|ftp)://[^\s:]+(?:\:\d+)?[^\s:]+)(:(?:title(?:=([^\]]+))?|barcode))?\]`),
			Handler: func(f *InlineFormatter, m []string) string { // [url:option]
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
			Pattern: regexp.MustCompile(`\[((?:https?|ftp):[^\s<>\]]+)\]`),
			Handler: func(f *InlineFormatter, m []string) string { // [url]
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
		{
			Pattern: regexp.MustCompile(`\[mailto:([^\s\@:?]+\@[^
\s\@:?]+(\?[^\s]+)?)\]`),
			Handler: func(f *InlineFormatter, m []string) string { // mailto
				uri := m[1]
				return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, uri, uri)
			},
		},
		{
			Pattern: regexp.MustCompile(`\[tex:([^\]]+)\]`),
			Handler: func(f *InlineFormatter, m []string) string { // tex
				tex := m[1]
				return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?cht=tx&chl=%s" alt="%s"/>`, url.QueryEscape(tex), html.EscapeString(tex))
			},
		},
		// 裸URL検出用パターンを追加
		{
			Pattern: regexp.MustCompile(`((?:https?|ftp):[^\s<>\"]+)`),
			Handler: func(f *InlineFormatter, m []string) string {
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
	}
}

func (f *InlineFormatter) Format(s string) string {
	s = strings.TrimPrefix(s, "\n")
	if len(f.rules) == 0 {
		f.rules = defaultInlineRules()
	}
	// すべてのパターンを | で結合した大きな正規表現を作成
	var patterns []string
	for _, r := range f.rules {
		patterns = append(patterns, r.Pattern.String())
	}
	bigRe := regexp.MustCompile(strings.Join(patterns, "|"))
	// 1パスで置換
	result := bigRe.ReplaceAllStringFunc(s, func(m string) string {
		for _, r := range f.rules {
			if sub := r.Pattern.FindStringSubmatch(m); sub != nil {
				return r.Handler(f, sub)
			}
		}
		return m
	})
	return result
}

func (f *InlineFormatter) Footnotes() []Footnote {
	return f.footnotes
}
