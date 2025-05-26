package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cho45/xatena-go/pkg/xatena"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: bench.go <inputfile>")
		os.Exit(1)
	}
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
		os.Exit(1)
	}
	x := xatena.NewXatena()

	ctx := context.Background()
	start := time.Now()
	var html string
	for i := 0; i < 1000; i++ {
		html = x.ToHTML(ctx, string(input))
	}
	elapsed := time.Since(start)

	fmt.Fprintf(os.Stderr, "go parse+format (1000x): %v ms\n", elapsed.Seconds())
	_ = html // 出力は不要なら省略
}
