package eventhandlers

import (
	"context"
	"goddd/internal/common/application"
	"goddd/internal/user/domain"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

type (
	UserCreatedHandler commonapplication.EventHandler[domain.UserCreatedEvent]
	userCreatedHandler struct{ fx.In }
)

func NewUserCreatedHandler(p userCreatedHandler) UserCreatedHandler {
	return &p
}

func (h *userCreatedHandler) Handle(
	ctx context.Context,
	log *slog.Logger,
	event domain.UserCreatedEvent,
) error {
	// Simulate work
	time.Sleep(1 * time.Second)
	log.InfoContext(ctx, "Stuff done!")

	return nil
}

type (
	UserCreatedHandler2 commonapplication.EventHandler[domain.UserCreatedEvent]
	userCreatedHandler2 struct{ fx.In }
)

func NewUserCreatedHandler2(p userCreatedHandler2) UserCreatedHandler2 {
	return &p
}

func (h *userCreatedHandler2) Handle(
	ctx context.Context,
	log *slog.Logger,
	event domain.UserCreatedEvent,
) error {
	// Simulate work
	time.Sleep(1 * time.Second)
	log.InfoContext(ctx, "Stuff done!")

	return nil
}
