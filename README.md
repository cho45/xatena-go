# xatena-go

Go実装の [はてな記法](https://help.hatenablog.com/entry/text-hatena-list ) → HTML変換ライブラリ & CLI

## 概要

- はてな記法をHTMLに変換するGoライブラリ・コマンドラインツールです。
- Perl版 [Text::Xatena](https://github.com/cho45/Text-Xatena ) のGo移植。
- パーサ・ノード構造・テストケースもPerl版に準拠。

## ディレクトリ構成

- `cmd/xatena-cli/` : CLIツール
- `internal/syntax/` : パーサ・ノード定義などコア実装
- `pkg/xatena/` : ライブラリAPI・テスト

## インストール

```sh
git clone https://github.com/cho45/xatena-go.git
cd xatena-go
go build ./cmd/xatena-cli
```

## 使い方

### CLI

```sh
cat sample.txt | ./xatena-cli
```

### ライブラリ

```go
import (
    "context"
    "github.com/cho45/xatena-go/pkg/xatena"
)

func main() {
    input := "* 見出し\n本文\n"
    x := xatena.NewXatena()
    html := x.ToHTML(context.Background(), input)
    // html を使う
}
```

インライン記法を追加する方法

```go
func getTitle(ctx context.Context, uri string) string {
    req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
    return ...
}

func main() {
    // タイトルハンドラ
    formatter := xatena.NewInlineFormatter(func(f *xatena.InlineFormatter) {
        f.SetTitleHandler(getTitle)
    })

    formatter.AddRule(InlineRule{
        Pattern: regexp.MustCompile(`\[custom:(.+?)\]`),
        Handler: func(ctx context.Context, f *InlineFormatter, m []string) string {
            return "<span>" + html.EscapeString(m[1]) + "</span>"
        },
    })

    x := xatena.NewXatenaWithInline(formatter)
    output := x.ToHTML(context.Background(), string(input))
}
```

はてな記法に挙動を近づける。(自動 p / br 挿入のルールが変化します)

```go
x := NewXatenaWithFields(NewInlineFormatter(), true)
...
```


## テスト

```sh
go test ./pkg/xatena
```

## 互換性・設計方針

- Perl版 Text::Xatena の仕様・挙動・テストにできるだけ準拠
- パーサは1行ずつ正規表現で消費するスタイル
- ノード構造・APIもPerl版に近い

## ライセンス

MIT License

---

## Design Doc: xatena-go の設計概要

### 全体設計

xatena-go は「はてな記法」を HTML へ変換するための Go 実装です。Perl版 Text::Xatena の設計・挙動を踏襲しつつ、Goらしい型安全性・拡張性・テスト容易性を重視しています。

### アーキテクチャ

- **Xatena 構造体**: パース・HTML変換の中心。BlockParser群やテンプレート、インライン整形器(Inline)を保持。
- **BlockParser インターフェース**: 各種ブロック要素(リスト、セクション、表など)ごとに実装。CanHandle/Parse メソッドで1行ずつパース。
- **Node/HasContent インターフェース**: パース結果のツリー構造を表現。各ノードは ToHTML で自身をHTML化。
- **Inline インターフェース**: インライン記法の整形器。Format メソッドでテキスト内のインライン要素をHTML化。

#### クラス図イメージ

```
Xatena
 ├─ Inline (interface)
 ├─ Templates (map)
 └─ BlockParsers ([]BlockParser)

BlockParser (interface)
 ├─ ListParser
 ├─ SectionParser
 ├─ ...

Node (interface)
 ├─ RootNode
 ├─ ListNode
 ├─ SectionNode
 ├─ ...

HasContent (interface)
 ├─ RootNode
 ├─ SectionNode
 ├─ ...
```

### パースの流れ

1. `Xatena#ToHTML(ctx, input)` でエントリ。
2. `parseXatena` で入力を正規化し、LineScanner で1行ずつ走査。
3. 各行ごとに BlockParser 群の CanHandle/Parse を順に適用。
4. マッチしない行は TextNode として追加。
5. Nodeツリー完成後、各ノードの ToHTML で再帰的にHTML化。
6. インライン要素は Inline.Format で整形。
7. テンプレートは Xatena.ExecuteTemplate で呼び出し。

### 拡張性

- BlockParser/Inline はインターフェース設計で差し替え・追加が容易。
- ブロックごとのテンプレートも map で管理し、用途に応じて差し替え可能。
- 新しい記法や出力形式の追加も最小限の実装で対応。

### テンプレートによる出力

- 各ブロック要素は html/template で出力を定義。
- Xatena.Templates でテンプレートを管理し、ExecuteTemplate で呼び出し。
- テンプレートのカスタマイズにより、HTML構造やクラス名の変更も容易。

### 互換性・Hatena互換モード

- `HatenaCompatible` フラグで、はてな記法の自動 p/br 挿入ルール等を切り替え。
- Perl版 Text::Xatena の仕様・テストに準拠。
- 既存のテストケースも移植し、互換性を担保。

---
