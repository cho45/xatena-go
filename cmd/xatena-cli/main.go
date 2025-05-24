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

var httpClient = &http.Client{Timeout: 5 * time.Second}

func getTitle(ctx context.Context, uri string) string {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[getTitle] failed to create request for %s: %v\n", uri, err)
		return uri
	}
	req.Header.Set("User-Agent", "xatena-cli/1.0 (+https://github.com/cho45/xatena-go)")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[getTitle] failed to GET %s: %v\n", uri, err)
		return uri // 失敗時はURLを返す
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "[getTitle] non-2xx status for %s: %d\n", uri, resp.StatusCode)
		return uri
	}
	// レスポンスボディのサイズ制限
	const maxBodySize = 2 * 1024 * 1024 // 2MB
	limitedBody := io.LimitReader(resp.Body, maxBodySize)
	z := html.NewTokenizer(limitedBody)
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
						return html.UnescapeString(title)
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
