package eventhandlers

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/post/domain"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

type (
	PostCreatedHandler commonapplication.EventHandler[domain.PostCreatedEvent]
	postCreatedHandler struct{ fx.In }
)

func NewPostCreatedHandler(p postCreatedHandler) PostCreatedHandler {
	return &p
}

func (h *postCreatedHandler) Handle(
	ctx context.Context,
	log *slog.Logger,
	event domain.PostCreatedEvent,
) error {
	// Simulate sending a notification
	time.Sleep(1 * time.Second)
	log.InfoContext(ctx, "Notification sent!", "post_id", event.GetAggregateId())

	return nil
}
