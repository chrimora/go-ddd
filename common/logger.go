package common

import (
	"context"
	"log/slog"
	"os"
)

const (
	RequestIdKey = "requestId"
	UserIdKey    = "userId"
)

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestId, ok := ctx.Value(RequestIdKey).(string); ok {
		r.AddAttrs(slog.String(RequestIdKey, requestId))
	}
	if userId, ok := ctx.Value(UserIdKey).(string); ok {
		r.AddAttrs(slog.String(UserIdKey, userId))
	}

	return h.Handler.Handle(ctx, r)
}

func NewLogger() *slog.Logger {
	base := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := &contextHandler{Handler: base}
	return slog.New(handler)
}
