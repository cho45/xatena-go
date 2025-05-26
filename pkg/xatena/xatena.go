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
	Inline           syntax.Inline
	Templates        map[string]*htmltpl.Template // テンプレート名→テンプレート
	HatenaCompatible bool                         // Hatena互換モードを使用するかどうか
	blockParsers     []syntax.BlockParser         // BlockParser のキャッシュ
}

func NewXatenaWithFields(inline syntax.Inline, hatenaCompatible bool) *Xatena {
	sectionTemplate := syntax.SectionTemplate
	if hatenaCompatible {
		sectionTemplate = syntax.HatenaCompatibleSectionTemplate
	}
	x := &Xatena{
		Inline: inline,
		Templates: map[string]*htmltpl.Template{
			"blockquote":     syntax.BlockquoteTemplate,
			"definitionlist": syntax.DefinitionListTemplate,
			"list":           syntax.ListTemplate,
			"section":        sectionTemplate,
			"table":          syntax.TableTemplate,
			"seemore":        syntax.SeeMoreTemplate,
			"pre":            syntax.PreTemplate,
			"stopp":          syntax.StopPTemplate,
			"superpre":       syntax.SuperPreTemplate,
			"comment":        syntax.CommentTemplate,
		},
		HatenaCompatible: hatenaCompatible,
	}
	x.blockParsers = []syntax.BlockParser{
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
	return x
}

func NewXatenaWithInline(inline syntax.Inline) *Xatena {
	return NewXatenaWithFields(inline, false)
}

func NewXatena() *Xatena {
	return NewXatenaWithInline(NewInlineFormatter())
}

// GetBlockParsers: Xatenaインスタンスを受け取る形に変更
func (x *Xatena) GetBlockParsers() []syntax.BlockParser {
	return x.blockParsers
}

// normalizeNewlines: \r\n, \r を \n に統一
func normalizeNewlines(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '\r' {
			if i+1 < len(s) && s[i+1] == '\n' {
				b.WriteByte('\n')
				i += 2
				continue
			}
			b.WriteByte('\n')
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// parseXatena: Xatenaインスタンスとcontext.Contextを受け取る
func (x *Xatena) parseXatena(ctx context.Context, input string) *syntax.RootNode {
	input = normalizeNewlines(input)
	parsers := x.GetBlockParsers()
	scanner := syntax.NewLineScanner(input)
	root := &syntax.RootNode{
		Content: make([]syntax.Node, 0, 4),
	}
	stack := []syntax.HasContent{root}
	for !scanner.EOF() {
		line := scanner.Peek()
		parent := stack[len(stack)-1]
		matched := false
		for _, parser := range parsers {
			if parser.CanHandle(line) {
				if parser.Parse(scanner, parent, &stack) {
					matched = true
					break
				}
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

func (x *Xatena) PreferHatenaCompatible() bool {
	return x.HatenaCompatible
}
