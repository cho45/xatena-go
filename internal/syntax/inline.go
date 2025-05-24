package syntax

import "context"

type Inline interface {
	Format(ctx context.Context, s string) string
}
