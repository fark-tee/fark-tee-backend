package context

import (
	"context"
	"log/slog"
)

// @WireSet("Infrastructure")
func NewContext(logger *slog.Logger) (context.Context, func()) {
	ctx := context.Background()

	return ctx, func() {
		ctx.Done()
	}
}
