package syntax

import (
	"context"
	"regexp"
	"strings"
)

type CallerOptions struct {
	stopp bool
}

type XatenaContext interface {
	GetInline() Inline
	ExecuteTemplate(name string, params map[string]interface{}) string
	PreferHatenaCompatible() bool
}

type Node interface {
	ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string
}

type HasContent interface {
	AddChild(n Node)
	GetContent() []Node
}

type RootNode struct {
	Content []Node
}

func SplitForBreak(text string) []string {
	if text == "" {
		return []string{}
	}
	parts := strings.Split(text, "\n")
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

var reToHTMLParagraph = regexp.MustCompile(`(\n{2,})`)
var reToHTMLParagraphLineBreak = regexp.MustCompile(`^\n+$`)

func ToHTMLParagraph(ctx context.Context, text string, xatena XatenaContext, options CallerOptions) string {
	text = xatena.GetInline().Format(ctx, text)
	if options.stopp {
		return text
	}
	parts := reSplitWithSep(reToHTMLParagraph, text)

	html := "<p>"
	for _, para := range parts {
		if reToHTMLParagraphLineBreak.MatchString(para) {
			html += "</p>" + strings.Repeat("<br />\n", (len(para)-2)) + "<p>"
		} else {
			html += strings.Join(SplitForBreak(para), "<br />\n")
		}
	}
	html += "</p>"
	return html
}

var reToHTMLParagraphHatenaCompatible = regexp.MustCompile(`(\n+)`)
var reToHTMLParagraphHatenaCompatibleLineBreak = regexp.MustCompile(`^(\n+)$`)

func ToHTMLParagraphHatenaCompatible(ctx context.Context, text string, xatena XatenaContext, options CallerOptions) string {
	text = xatena.GetInline().Format(ctx, text)
	text = strings.TrimSuffix(text, "\n") // Remove trailing newline
	if options.stopp {
		return text
	}
	parts := reSplitWithSep(reToHTMLParagraphHatenaCompatible, text)

	html := "<p>"
	for _, para := range parts {
		if m := reToHTMLParagraphHatenaCompatibleLineBreak.FindStringSubmatch(para); m != nil {
			brs := len(m[1]) - 2
			if brs < 0 {
				brs = 0
			}
			html += "</p>" + strings.Repeat("<br />\n", brs) + "<p>"
		} else {
			html += para
		}
	}
	html += "</p>"
	return html
}

//	local *Text::Xatena::Node::as_html_paragraph = sub {
//	    my ($self, $context, $text, %opts) = @_;
//	    $text = $context->inline->format($text, %opts);
//
//	    $text =~ s{\n$}{}g;
//	    if ($opts{stopp}) {
//	        $text;
//	    } else {
//	        "<p>" . join("",
//	            map {
//	                if (/^(\n+)$/) {
//	                    "</p>" . ("<br />\n" x (length($1) - 2)) . "<p>";
//	                } else {
//	                    $_;
//	                }
//	            }
//	            split(/(\n+)/, $text)
//	        ) . "</p>\n";
//	    }
//	};
func ContentToHTML(r HasContent, ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	html := ""
	var textBuf []string
	flushParagraph := func() {
		hasText := strings.Join(textBuf, "") != ""
		if hasText {
			if xatena.PreferHatenaCompatible() {
				html += ToHTMLParagraphHatenaCompatible(ctx, strings.Join(textBuf, "\n"), xatena, options)
			} else {
				html += ToHTMLParagraph(ctx, strings.Join(textBuf, "\n"), xatena, options)
			}
		}
		textBuf = nil
	}
	for _, n := range r.GetContent() {
		if t, ok := n.(*TextNode); ok {
			textBuf = append(textBuf, t.Text)
		} else {
			flushParagraph()
			html += n.ToHTML(ctx, xatena, options)
		}
	}
	flushParagraph()
	return html
}

// reSplitWithSep splits s by the regexp re, including the separator (match) in the result slice.
// e.g. reSplitWithSep(regexp.MustCompile(`(\n{2,})`), "a\n\nb\n\n\nc") => ["a", "\n\n", "b", "\n\n\n", "c"]
func reSplitWithSep(re *regexp.Regexp, s string) []string {
	result := []string{}
	last := 0
	indices := re.FindAllStringIndex(s, -1)
	for _, idx := range indices {
		if last < idx[0] {
			result = append(result, s[last:idx[0]])
		}
		result = append(result, s[idx[0]:idx[1]])
		last = idx[1]
	}
	if last < len(s) {
		result = append(result, s[last:])
	}
	return result
}

type BlockParser interface {
	Parse(scanner *LineScanner, parent HasContent, stack *[]HasContent) bool
}

func (r *RootNode) AddChild(n Node) {
	r.Content = append(r.Content, n)
}

func (r *RootNode) GetContent() []Node {
	return r.Content
}

func (r *RootNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	return ContentToHTML(r, ctx, xatena, options)
}

type TextNode struct {
	Text string
}

func (t *TextNode) ToHTML(ctx context.Context, xatena XatenaContext, options CallerOptions) string {
	panic("TextNode does not support ToHTML")
}
