package commonapplication

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Internal
type command[I any] interface {
	Handle(ctx context.Context, log *slog.Logger, input I) (uuid.UUID, error)
}

// Wrapped
type CommandI[I any] interface {
	Handle(ctx context.Context, input I) (uuid.UUID, error)
}

type commandWrapper[I any] struct {
	log     *slog.Logger
	handler command[I]
}

func (c *commandWrapper[I]) name() string {
	return reflect.TypeOf(c.handler).Elem().Name()
}
func (c *commandWrapper[I]) Handle(ctx context.Context, input I) (uuid.UUID, error) {
	log := c.log.With("command", c.name())
	log.InfoContext(ctx, "Command start")

	start := time.Now()
	id, err := c.handler.Handle(ctx, log, input)
	if err != nil {
		log.ErrorContext(ctx, "Command error", "err", err)
	} else {
		log.InfoContext(ctx, "Command complete", "duration", time.Since(start))
	}

	return id, err
}

func NewCommand[I any](log *slog.Logger, h command[I]) *commandWrapper[I] {
	return &commandWrapper[I]{log: log, handler: h}
}
