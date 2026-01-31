package application

import (
	"context"
	"goddd/internal/common/infrastructure"
	"goddd/internal/outbox/domain"

	"github.com/ThreeDotsLabs/watermill/message"
)

type EventPublisherI interface {
	Publish(context.Context, *domain.OutboxEvent) error
}

type EventPublisher EventPublisherI
type eventPublisher struct {
	publisher message.Publisher
}

func NewEventPublisher(pub message.Publisher) EventPublisher {
	return &eventPublisher{
		publisher: pub,
	}
}

func (p *eventPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	msg := message.NewMessage(event.ID.String(), event.Payload)

	trace, err := commoninfrastructure.NewTraceCtxFromJson(event.EventContext)
	if err != nil {
		return err
	}
	msg.Metadata.Set(commoninfrastructure.RequestIdKey, trace.RequestId)
	msg.Metadata.Set(commoninfrastructure.UserIdKey, trace.UserId)

	return p.publisher.Publish(event.EventType, msg)
}
