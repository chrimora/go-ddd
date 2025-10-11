package eventhandlers

import (
	"context"
	"goddd/internal/user/domain"
	"log/slog"
	"time"
)

type UserCreatedHandler struct{}

func NewUserCreatedHandler() *UserCreatedHandler {
	return &UserCreatedHandler{}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, log *slog.Logger, event domain.UserCreatedEvent) error {
	// Simulate work
	log.InfoContext(ctx, "Doing stuff!")
	time.Sleep(1 * time.Second)

	return nil
}

type UserCreatedHandler2 struct{}

func NewUserCreatedHandler2() *UserCreatedHandler2 {
	return &UserCreatedHandler2{}
}

func (h *UserCreatedHandler2) Handle(ctx context.Context, log *slog.Logger, event domain.UserCreatedEvent) error {
	// Simulate work
	log.InfoContext(ctx, "Doing stuff!")
	time.Sleep(1 * time.Second)

	return nil
}
