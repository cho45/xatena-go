package xatena

import (
	"context"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type InlineRule struct {
	Pattern *regexp.Regexp
	Handler func(ctx context.Context, f *InlineFormatter, m []string) string
}

type InlineFormatter struct {
	footnotes    []Footnote
	rules        []InlineRule
	bigRe        *regexp.Regexp
	titleHandler func(ctx context.Context, uri string) string
}

type Footnote struct {
	Number int
	Note   string
	Title  string
}

func (f *InlineFormatter) AddRule(rule InlineRule) {
	f.rules = append(f.rules, rule)
	f.bigRe = nil // Reset the big regex cache
}

func (f *InlineFormatter) AddRuleAt(index int, rule InlineRule) {
	if index < 0 || index > len(f.rules) {
		index = len(f.rules)
	}
	f.rules = append(f.rules[:index], append([]InlineRule{rule}, f.rules[index:]...)...)
	f.bigRe = nil // Reset the big regex cache
}

func defaultTitleHandler(ctx context.Context, uri string) string {
	return uri
}

func NewInlineFormatter(opts ...func(*InlineFormatter)) *InlineFormatter {
	f := &InlineFormatter{
		footnotes:    []Footnote{},
		titleHandler: defaultTitleHandler,
	}
	f.rules = defaultInlineRules(f)
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func defaultInlineRules(f *InlineFormatter) []InlineRule {
	return []InlineRule{
		{
			Pattern: regexp.MustCompile(`\[\]([\s\S]*?)\[\]`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return m[1] },
		},
		{
			Pattern: regexp.MustCompile(`\(\(\(.*?\)\)\)`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return m[0][1 : len(m[0])-1] },
		},
		{
			Pattern: regexp.MustCompile(`\)\(\(.*?\)\)\(`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return m[0][1 : len(m[0])-1] },
		},
		{
			Pattern: regexp.MustCompile(`\(\((.+?)\)\)`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				note := m[1]
				title := html.EscapeString(note)
				f.footnotes = append(f.footnotes, Footnote{Number: len(f.footnotes) + 1, Note: note, Title: title})
				return fmt.Sprintf(`<a href="#fn%d" title="%s">*%d</a>`, len(f.footnotes), html.EscapeString(title), len(f.footnotes))
			},
		},
		{
			Pattern: regexp.MustCompile(`(?i)<a[^>]+>[\s\S]*?</a>`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return m[0] },
		},
		{
			Pattern: regexp.MustCompile(`<!--.*?-->`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return "<!-- -->" },
		},
		{
			Pattern: regexp.MustCompile(`(?i)<[^>]+>`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string { return m[0] },
		},
		{
			Pattern: regexp.MustCompile(`\[((?:https?|ftp)://[^\s:]+(?:\:\d+)?[^\s:]+)(:(?:title(?:=([^\]]+))?|barcode))?\]`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				uri, opt, title := m[1], m[2], m[3]
				if opt == ":barcode" {
					return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?chs=150x150&cht=qr&chl=%s" title="%s"/>`, url.QueryEscape(uri), html.EscapeString(uri))
				}
				if strings.HasPrefix(opt, ":title") {
					if title != "" {
						return fmt.Sprintf(`<a href="%s">%s</a>`, uri, html.EscapeString(title))
					}
					return fmt.Sprintf(`<a href="%s">%s</a>`, uri, html.EscapeString(f.titleHandler(ctx, uri)))
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, html.EscapeString(uri))
			},
		},
		{
			Pattern: regexp.MustCompile(`\[((?:https?|ftp):[^\s<>\]]+)\]`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
		{
			Pattern: regexp.MustCompile(`\[mailto:([^\s\@:?]+\@[^\s\@:?]+(\?[^\s]+)?)\]`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				uri := m[1]
				return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, uri, uri)
			},
		},
		{
			Pattern: regexp.MustCompile(`\[tex:([^\]]+)\]`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				tex := m[1]
				return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?cht=tx&chl=%s" alt="%s"/>`, url.QueryEscape(tex), html.EscapeString(tex))
			},
		},
		{
			Pattern: regexp.MustCompile(`((?:https?|ftp):[^\s<>\"]+)`),
			Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
				uri := m[1]
				if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
					return m[0]
				}
				return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
			},
		},
	}
}

func (f *InlineFormatter) Format(ctx context.Context, s string) string {
	s = strings.TrimPrefix(s, "\n")
	if len(f.rules) == 0 {
		f.rules = defaultInlineRules(f)
	}
	var patterns []string
	for _, r := range f.rules {
		patterns = append(patterns, r.Pattern.String())
	}
	if f.bigRe == nil {
		f.bigRe = regexp.MustCompile(strings.Join(patterns, "|"))
	}
	result := f.bigRe.ReplaceAllStringFunc(s, func(m string) string {
		for _, r := range f.rules {
			if sub := r.Pattern.FindStringSubmatch(m); sub != nil {
				return r.Handler(ctx, f, sub)
			}
		}
		return m
	})
	return result
}

func (f *InlineFormatter) Footnotes() []Footnote {
	return f.footnotes
}

func (f *InlineFormatter) SetTitleHandler(handler func(ctx context.Context, uri string) string) {
	f.titleHandler = handler
}
