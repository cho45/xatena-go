package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/cho45/xatena-go/pkg/xatena"
)

func getTitle(uri string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[getTitle] failed to GET %s: %v\n", uri, err)
		return uri // 失敗時はURLを返す
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "[getTitle] non-2xx status for %s: %d\n", uri, resp.StatusCode)
		return uri
	}
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "[getTitle] HTML parse error for %s: %v\n", uri, err)
			}
			return uri
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "title" {
				if z.Next() == html.TextToken {
					title := strings.TrimSpace(z.Token().Data)
					if title != "" {
						return title
					}
				}
			}
		}
	}
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read stdin: %v\n", err)
		os.Exit(1)
	}

	formatter := xatena.NewInlineFormatter(func(f *xatena.InlineFormatter) {
		f.SetTitleHandler(getTitle)
	})

	x := xatena.NewXatenaWithInline(formatter)
	output := x.ToHTML(context.Background(), string(input))
	fmt.Print(strings.TrimRight(output, "\n"))
}
