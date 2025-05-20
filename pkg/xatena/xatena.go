package xatena

import (
	"github.com/cho45/xatena-go/internal/syntax"
)

// GetBlockParsers returns all block-level parsers in Xatena order
func GetBlockParsers() []syntax.BlockParser {
	return []syntax.BlockParser{
		// &syntax.SeeMore{},
		&syntax.SuperPreParser{},
		&syntax.StopPParser{},
		&syntax.BlockquoteParser{},
		&syntax.PreParser{},
		&syntax.ListParser{},
		&syntax.DefinitionListParser{},
		&syntax.TableParser{},
		&syntax.SectionParser{},
		// &syntax.CommentParser{},
	}
}

func parseXatena(input string) *syntax.RootNode {
	parsers := GetBlockParsers()
	scanner := syntax.NewLineScanner(input)
	root := &syntax.RootNode{}
	stack := []syntax.BlockNode{root}
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

func Format(input string) string {
	node := parseXatena(input)
	return node.ToHTML()
}
