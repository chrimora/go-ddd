package eventhandlers

import (
	"context"
	"goddd/internal/domain/user"
	"log/slog"
	"time"
)

type UserCreatedHandler struct {
	log *slog.Logger
}

func NewUserCreatedHandler(log *slog.Logger) *UserCreatedHandler {
	return &UserCreatedHandler{log: log}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, event user.UserCreatedEvent) error {
	// Simulate work
	time.Sleep(1 * time.Second)
	h.log.InfoContext(ctx, "Doing stuff!")

	return nil
}

type UserCreatedHandler2 struct {
	log *slog.Logger
}

func NewUserCreatedHandler2(log *slog.Logger) *UserCreatedHandler2 {
	return &UserCreatedHandler2{log: log}
}

func (h *UserCreatedHandler2) Handle(ctx context.Context, event user.UserCreatedEvent) error {
	// Simulate work
	time.Sleep(1 * time.Second)
	h.log.InfoContext(ctx, "Doing stuff!")

	return nil
}
