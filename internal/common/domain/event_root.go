package commondomain

import (
	"github.com/google/uuid"
)

type EventType string
type DomainEventI interface {
	GetEventType() EventType
	GetAggregateType() string
	GetAggregateId() uuid.UUID
}

type DomainEventRoot struct {
	AggregateId uuid.UUID `json:"aggregate_id"`
}

func (d DomainEventRoot) GetAggregateId() uuid.UUID {
	return d.AggregateId
}
