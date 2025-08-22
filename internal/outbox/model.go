package outbox

import (
	"context"
	"encoding/json"
	"gotemplate/internal/common"
	outboxdb "gotemplate/internal/outbox/db"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

type EventType string
type DomainEventI interface {
	Type() EventType
}

type OutboxEvent struct {
	ID          int
	EventType   EventType
	Payload     []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Retries     int32
	ProcessedAt *time.Time
}

func fromDB(outbox outboxdb.EventOutbox) *OutboxEvent {
	return &OutboxEvent{
		ID:          int(outbox.ID),
		EventType:   EventType(outbox.EventType),
		Payload:     outbox.Payload,
		CreatedAt:   outbox.CreatedAt,
		UpdatedAt:   outbox.UpdatedAt,
		Retries:     outbox.Retries,
		ProcessedAt: outbox.ProcessedAt,
	}
}

func CreateEventPayload(ctx context.Context, payloadStruct any) ([]byte, error) {
	var payloadMap map[string]any
	err := mapstructure.Decode(payloadStruct, &payloadMap)
	if err != nil {
		return nil, err
	}

	// TODO
	userID, _ := ctx.Value(common.UserIdKey).(string)
	requestID, _ := ctx.Value(common.RequestIdKey).(string)

	payloadMap[common.UserIdKey] = userID
	payloadMap[common.RequestIdKey] = requestID

	return json.Marshal(payloadMap)
}
