package xatena

import (
	"context"

	"github.com/cho45/xatena-go/internal/syntax"
)

// Xatena 構造体: InlineFormatter などを保持
// 今後オプションや拡張もここに集約

type Xatena struct {
	Inline *syntax.InlineFormatter
}

func NewXatena() *Xatena {
	return &Xatena{
		Inline: syntax.NewInlineFormatter(),
	}
}

// GetBlockParsers: Xatenaインスタンスを受け取る形に変更
func (x *Xatena) GetBlockParsers() []syntax.BlockParser {
	return []syntax.BlockParser{
		&syntax.SuperPreParser{},
		&syntax.StopPParser{},
		&syntax.BlockquoteParser{},
		&syntax.PreParser{},
		&syntax.ListParser{},
		&syntax.DefinitionListParser{},
		&syntax.TableParser{},
		&syntax.SectionParser{},
	}
}

// parseXatena: Xatenaインスタンスとcontext.Contextを受け取る
func (x *Xatena) parseXatena(ctx context.Context, input string) *syntax.RootNode {
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
	return node.ToHTML(ctx)
}
