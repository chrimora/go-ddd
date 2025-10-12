package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateID   uuid.UUID
	AggregateType string
	EventContext  []byte
	EventType     string
	Payload       []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (e *OutboxEvent) String() string {
	return fmt.Sprintf("OutboxEvent[id: %s, type: %s]", e.ID, e.EventType)
}
