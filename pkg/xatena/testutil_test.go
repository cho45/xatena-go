package xatena

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// テストフィクスチャのパース用データ型
// 例: === testname\n--- input\n...\n--- expected\n...
type TestBlock struct {
	Name     string
	Sections map[string]string // section名→内容
}

// === [name]\n--- [section]\n...\n--- [section]\n... 形式のテストフィクスチャをパース
func parseTestBlocks(data string) []TestBlock {
	var blocks []TestBlock
	parts := strings.Split(data, "=== ")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		lines := strings.Split(part, "\n")
		name := lines[0]
		sections := make(map[string]string)
		var currentSection string
		var buf []string
		for _, line := range lines[1:] {
			if strings.HasPrefix(line, "--- ") {
				if currentSection != "" {
					sections[currentSection] = strings.TrimRight(strings.Join(buf, "\n"), "\n")
				}
				currentSection = strings.TrimSpace(line[4:])
				buf = buf[:0]
			} else {
				buf = append(buf, line)
			}
		}
		if currentSection != "" {
			sections[currentSection] = strings.TrimRight(strings.Join(buf, "\n"), "\n")
		}
		blocks = append(blocks, TestBlock{Name: name, Sections: sections})
	}
	return blocks
}

// EqualHTML は2つのHTML文字列を構造的に比較し、異なればt.Errorfで詳細を出力する
func EqualHTML(t *testing.T, a, b string) {
	nodeA, errA := html.Parse(strings.NewReader(a))
	nodeB, errB := html.Parse(strings.NewReader(b))
	if errA != nil || errB != nil {
		t.Errorf("HTML parse error: %v / %v", errA, errB)
		return
	}
	if !equalHTMLNode(nodeA, nodeB) {
		t.Errorf("HTML not equal (structure):\nA: %s\nB: %s", a, b)
	}
}

func equalHTMLNode(a, b *html.Node) bool {
	if a.Type != b.Type || a.Data != b.Data {
		return false
	}
	// 属性比較
	if len(a.Attr) != len(b.Attr) {
		return false
	}
	for i := range a.Attr {
		if a.Attr[i].Key != b.Attr[i].Key || a.Attr[i].Val != b.Attr[i].Val {
			return false
		}
	}
	// 子ノード比較（空白だけのテキストノードは無視）
	nextNonEmptyText := func(n *html.Node) *html.Node {
		for n != nil && n.Type == html.TextNode && strings.TrimSpace(n.Data) == "" {
			n = n.NextSibling
		}
		return n
	}
	ca, cb := a.FirstChild, b.FirstChild
	for {
		ca = nextNonEmptyText(ca)
		cb = nextNonEmptyText(cb)
		if ca == nil || cb == nil {
			break
		}
		if ca.Type == html.TextNode && cb.Type == html.TextNode {
			if strings.TrimSpace(ca.Data) != strings.TrimSpace(cb.Data) {
				return false
			}
		} else {
			if !equalHTMLNode(ca, cb) {
				return false
			}
		}
		ca = ca.NextSibling
		cb = cb.NextSibling
	}
	ca = nextNonEmptyText(ca)
	cb = nextNonEmptyText(cb)
	if ca != nil || cb != nil {
		return false
	}
	return true
}
