package commoninfrastructure

import (
	"context"
	"goddd/internal/config"
	"log/slog"
	"os"
)

type contextHandler struct {
	slog.Handler
	service *config.ServiceConfig
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("service", h.service.Name))

	if rc, ok := ctx.Value(RequestContextKey).(RequestContext); ok {
		r.AddAttrs(slog.String(requestIdMetaKey, rc.RequestId.String()))
		r.AddAttrs(slog.String(userIdMetaKey, rc.UserId.String()))
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

func NewLogger(service *config.ServiceConfig) *slog.Logger {
	var base slog.Handler
	if service.IsProd() {
		base = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true})
	} else {
		base = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	handler := &contextHandler{Handler: base, service: service}
	return slog.New(handler)
}
