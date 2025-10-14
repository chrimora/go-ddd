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
	UserCreatedHandlerParams struct{ fx.In }
	UserCreatedHandler       commonapplication.EventHandler[domain.UserCreatedEvent]
	userCreatedHandler       struct{}
)

func NewUserCreatedHandler(p UserCreatedHandlerParams) UserCreatedHandler {
	return &userCreatedHandler{}
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
	UserCreatedHandlerParams2 struct{ fx.In }
	UserCreatedHandler2       commonapplication.EventHandler[domain.UserCreatedEvent]
	userCreatedHandler2       struct{}
)

func NewUserCreatedHandler2(p UserCreatedHandlerParams2) UserCreatedHandler2 {
	return &userCreatedHandler2{}
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
