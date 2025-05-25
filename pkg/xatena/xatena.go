package xatena

import (
	"context"
	htmltpl "html/template"
	"strings"

	"github.com/cho45/xatena-go/internal/syntax"
)

// Xatena 構造体: InlineFormatter などを保持
// 今後オプションや拡張もここに集約

type Xatena struct {
	Inline    syntax.Inline
	Templates map[string]*htmltpl.Template // テンプレート名→テンプレート
}

func NewXatenaWithInline(inline syntax.Inline) *Xatena {
	return &Xatena{
		Inline: inline,
		Templates: map[string]*htmltpl.Template{
			"blockquote":     syntax.BlockquoteTemplate,
			"definitionlist": syntax.DefinitionListTemplate,
			"list":           syntax.ListTemplate,
			"section":        syntax.SectionTemplate,
			"table":          syntax.TableTemplate,
			"seemore":        syntax.SeeMoreTemplate,
			"pre":            syntax.PreTemplate,
			"stopp":          syntax.StopPTemplate,
			"superpre":       syntax.SuperPreTemplate,
			"comment":        syntax.CommentTemplate,
		},
	}
}

func NewXatena() *Xatena {
	return NewXatenaWithInline(NewInlineFormatter())
}

// GetBlockParsers: Xatenaインスタンスを受け取る形に変更
func (x *Xatena) GetBlockParsers() []syntax.BlockParser {
	return []syntax.BlockParser{
		&syntax.SeeMoreParser{},
		&syntax.SuperPreParser{},
		&syntax.StopPParser{},
		&syntax.BlockquoteParser{},
		&syntax.PreParser{},
		&syntax.ListParser{},
		&syntax.DefinitionListParser{},
		&syntax.TableParser{},
		&syntax.SectionParser{},
		&syntax.CommentParser{},
	}
}

// parseXatena: Xatenaインスタンスとcontext.Contextを受け取る
func (x *Xatena) parseXatena(ctx context.Context, input string) *syntax.RootNode {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	parsers := x.GetBlockParsers()
	scanner := syntax.NewLineScanner(input)
	root := &syntax.RootNode{}
	stack := []syntax.HasContent{root}
	for !scanner.EOF() {
		parent := stack[len(stack)-1]
		matched := false
		for _, parser := range parsers {
			if parser.Parse(scanner, parent, &stack) {
				matched = true
				break
			}
		}
		if !matched {
			parent.AddChild(&syntax.TextNode{Text: scanner.Next()})
		}
	}
	return root
}

// ToHTML: Xatenaインスタンスとcontext.Contextを渡す
func (x *Xatena) ToHTML(ctx context.Context, input string) string {
	node := x.parseXatena(ctx, input)
	return node.ToHTML(ctx, x, syntax.CallerOptions{})
}

func (x *Xatena) GetInline() syntax.Inline {
	return x.Inline
}

func (x *Xatena) ExecuteTemplate(name string, params map[string]interface{}) string {
	tmpl, ok := x.Templates[name]
	if !ok {
		return `<div class="xatena-template-error">template not found: ` + htmltpl.HTMLEscapeString(name) + `</div>`
	}
	var sb strings.Builder
	err := tmpl.Execute(&sb, params)
	if err != nil {
		return `<div class="xatena-template-error">template error: ` + htmltpl.HTMLEscapeString(err.Error()) + `</div>`
	}
	return sb.String()
}
