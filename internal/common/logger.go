package common

import (
	"context"
	"gotemplate/internal/config"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
	service config.ServiceConfig
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("service", h.service.Name))

	serviceCtx, err := NewServiceCtx(ctx)
	if err != nil {
		return err
	}

	r.AddAttrs(slog.String(RequestIdKey, serviceCtx.RequestId))
	// r.AddAttrs(slog.String(UserIdKey, serviceCtx.UserId))

	return h.Handler.Handle(ctx, r)
}

func NewLogger(service config.ServiceConfig) *slog.Logger {
	// Prod
	// base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	base := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler := &contextHandler{Handler: base, service: service}
	return slog.New(handler)
}
