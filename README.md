# xatena-go

Go実装の [はてな記法](https://hatena.github.io/text-hatena/) → HTML変換ライブラリ & CLI

## 概要

- はてな記法をHTMLに変換するGoライブラリ・コマンドラインツールです。
- Perl版 [Text::Xatena](https://github.com/hatena/Text-Xatena) のGo移植。
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
import "github.com/cho45/xatena-go/pkg/xatena"

html := xatena.Render("* 見出し\n本文\n")
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
