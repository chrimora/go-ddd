package outbox

import (
	"gotemplate/internal/domain/common"
	outboxdb "gotemplate/internal/outbox/db"
	"time"
)

type OutboxEvent struct {
	ID          int
	EventType   common.EventType
	Payload     []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Retries     int32
	ProcessedAt *time.Time
}

func fromDB(outbox outboxdb.EventOutbox) *OutboxEvent {
	return &OutboxEvent{
		ID:          int(outbox.ID),
		EventType:   common.EventType(outbox.EventType),
		Payload:     outbox.Payload,
		CreatedAt:   outbox.CreatedAt,
		UpdatedAt:   outbox.UpdatedAt,
		Retries:     outbox.Retries,
		ProcessedAt: outbox.ProcessedAt,
	}
}
