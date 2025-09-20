package common

import (
	"github.com/google/uuid"
)

type EventType string
type DomainEventI interface {
	GetType() EventType
	GetAggregateId() uuid.UUID
}
type DomainEventRoot struct {
	EventType   EventType
	AggregateId uuid.UUID
}

func (d DomainEventRoot) GetType() EventType {
	return d.EventType
}
func (d DomainEventRoot) GetAggregateId() uuid.UUID {
	return d.AggregateId
}
