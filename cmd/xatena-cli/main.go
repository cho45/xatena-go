package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cho45/xatena-go/pkg/xatena"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v\n", err)
		os.Exit(1)
	}

	x := xatena.NewXatena()
	output := x.ToHTML(context.Background(), string(input))
	fmt.Print(strings.TrimRight(output, "\n"))
}
