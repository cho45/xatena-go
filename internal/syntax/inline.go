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
	// Perlのmatch順に正規表現を適用
	// 1. [[]...[]] などの特殊なアンリンク
	s = regexp.MustCompile(`\[\]([\s\S]*?)\[\]`).ReplaceAllStringFunc(s, func(m string) string {
		return m[2 : len(m)-2]
	})
	s = regexp.MustCompile(`\(\(\(.*?\)\)\)`).ReplaceAllStringFunc(s, func(m string) string {
		return m[1 : len(m)-1]
	})
	s = regexp.MustCompile(`\)\(\(.*?\)\)\(`).ReplaceAllStringFunc(s, func(m string) string {
		return m[1 : len(m)-1]
	})
	// 2. footnote
	s = regexp.MustCompile(`\(\((.+?)\)\)`).ReplaceAllStringFunc(s, func(m string) string {
		note := m[2 : len(m)-2]
		title := stripTags(note)
		f.footnotes = append(f.footnotes, Footnote{Number: len(f.footnotes) + 1, Note: note, Title: title})
		return fmt.Sprintf(`<a href="#fn%d" title="%s">*%d</a>`, len(f.footnotes), html.EscapeString(title), len(f.footnotes))
	})
	// 3. <a...>...</a> そのまま
	s = regexp.MustCompile(`(?i)<a[^>]+>[\s\S]*?</a>`).ReplaceAllStringFunc(s, func(m string) string {
		return m
	})
	// 4. <!-- ... -->
	s = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(s, "<!-- -->")
	// 5. <...> そのまま
	s = regexp.MustCompile(`(?i)<[^>]+>`).ReplaceAllStringFunc(s, func(m string) string {
		return m
	})
	// 6. [url:option] (barcode, title)
	s = regexp.MustCompile(`\[((?:https?|ftp)://[^\s:]+(?:\:\d+)?[^\s:]+)(:(?:title(?:=([^\]]+))?|barcode))?\]`).ReplaceAllStringFunc(s, func(m string) string {
		if !strings.HasPrefix(m, "[") || !strings.HasSuffix(m, "]") {
			return m
		}
		re := regexp.MustCompile(`\[((?:https?|ftp)://[^\s:]+(?:\:\d+)?[^\s:]+)(:(?:title(?:=([^\]]+))?|barcode))?\]`)
		m2 := re.FindStringSubmatch(m)
		if m2 == nil {
			return m
		}
		uri, opt, title := m2[1], m2[2], m2[3]
		if opt == ":barcode" {
			return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?chs=150x150&cht=qr&chl=%s" title="%s"/>`, url.QueryEscape(uri), html.EscapeString(uri))
		}
		if strings.HasPrefix(opt, ":title") {
			if title == "" && f.aggressive {
				// aggressive title fetch (未実装)
			}
			return fmt.Sprintf(`<a href="%s">%s</a>`, uri, html.EscapeString(title))
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
	})
	// 7. [url] (通常リンク)
	s = regexp.MustCompile(`\[((?:https?|ftp):[^\s<>\]]+)\]`).ReplaceAllStringFunc(s, func(m string) string {
		if !strings.HasPrefix(m, "[") || !strings.HasSuffix(m, "]") {
			return m
		}
		re := regexp.MustCompile(`\[((?:https?|ftp):[^\s<>\]]+)\]`)
		m2 := re.FindStringSubmatch(m)
		if m2 == nil {
			return m
		}
		uri := m2[1]
		if strings.HasSuffix(uri, ":barcode") || strings.HasPrefix(uri, ":title") {
			return m
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, uri, uri)
	})
	// 8. mailto
	s = regexp.MustCompile(`\[mailto:([^\s\@:?]+\@[^\s\@:?]+(\?[^\s]+)?)\]`).ReplaceAllStringFunc(s, func(m string) string {
		if !strings.HasPrefix(m, "[") || !strings.HasSuffix(m, "]") {
			return m
		}
		re := regexp.MustCompile(`\[mailto:([^\s\@:?]+\@[^\s\@:?]+(\?[^\s]+)?)\]`)
		m2 := re.FindStringSubmatch(m)
		if m2 == nil {
			return m
		}
		uri := m2[1]
		return fmt.Sprintf(`<a href="mailto:%s">%s</a>`, uri, uri)
	})
	// 9. tex
	s = regexp.MustCompile(`\[tex:([^\]]+)\]`).ReplaceAllStringFunc(s, func(m string) string {
		re := regexp.MustCompile(`\[tex:([^\]]+)\]`)
		m2 := re.FindStringSubmatch(m)
		if m2 == nil {
			return m
		}
		tex := m2[1]
		return fmt.Sprintf(`<img src="http://chart.apis.google.com/chart?cht=tx&chl=%s" alt="%s"/>`, url.QueryEscape(tex), html.EscapeString(tex))
	})
	return s
}

func stripTags(s string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(s, "")
}

func (f *InlineFormatter) Footnotes() []Footnote {
	return f.footnotes
}
