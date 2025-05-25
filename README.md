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
