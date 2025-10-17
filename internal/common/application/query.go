package commonapplication

import (
	"context"
	"log/slog"
	"reflect"
	"time"
)

// Internal
type query[I any, O any] interface {
	Handle(ctx context.Context, log *slog.Logger, input I) (O, error)
}

// Wrapped
type QueryI[I any, O any] interface {
	Handle(ctx context.Context, input I) (O, error)
}

type queryWrapper[I any, O any] struct {
	log     *slog.Logger
	handler query[I, O]
}

func (c *queryWrapper[I, O]) name() string {
	return reflect.TypeOf(c.handler).Elem().Name()
}
func (c *queryWrapper[I, O]) Handle(ctx context.Context, input I) (O, error) {
	log := c.log.With("query", c.name())
	log.InfoContext(ctx, "Query start")

	start := time.Now()
	res, err := c.handler.Handle(ctx, log, input)
	if err != nil {
		log.ErrorContext(ctx, "Query error", "err", err)
	} else {
		log.InfoContext(ctx, "Query complete", "duration", time.Since(start))
	}

	return res, err
}

func NewQuery[I any, O any](log *slog.Logger, h query[I, O]) *queryWrapper[I, O] {
	return &queryWrapper[I, O]{log: log, handler: h}
}
