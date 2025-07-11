//go:build js && wasm
// +build js,wasm

package main

import (
	"context"
	"syscall/js"

	"github.com/cho45/xatena-go/pkg/xatena"
)

var x *xatena.Xatena

func init() {
	x = xatena.NewXatena()
}

func toHTMLWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "error: need 1 argument (input string)"
	}
	input := args[0].String()
	html := x.ToHTML(context.Background(), input)
	return html
}

func main() {
	js.Global().Set("xatenaToHTML", js.FuncOf(toHTMLWrapper))
	select {} // keep running
}
