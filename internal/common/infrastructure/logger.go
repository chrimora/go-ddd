package infrastructure

import (
	"context"
	"goddd/internal/config"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
	service config.ServiceConfig
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("service", h.service.Name))

	traceCtx := NewTraceCtx(ctx)

	if traceCtx.RequestId != "" {
		r.AddAttrs(slog.String(RequestIdKey, traceCtx.RequestId))
	}
	if traceCtx.UserId != "" {
		r.AddAttrs(slog.String(UserIdKey, traceCtx.UserId))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{
		Handler: h.Handler.WithAttrs(attrs),
		service: h.service,
	}
}
func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{
		Handler: h.Handler.WithGroup(name),
		service: h.service,
	}
}

func NewLogger(service config.ServiceConfig) *slog.Logger {
	// Prod
	// base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	base := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := &contextHandler{Handler: base, service: service}
	return slog.New(handler)
}
