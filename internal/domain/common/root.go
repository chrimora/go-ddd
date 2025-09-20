package common

import (
	"context"
	"gotemplate/internal/common"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OutboxRepositoryI interface {
	WithTx(pgx.Tx) OutboxRepositoryI
	Create(context.Context, common.ServiceContext, DomainEventI) error
}

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

type AggregateRootI interface {
	PullEvents() []DomainEventI
}
type AggregateRoot struct {
	ID     uuid.UUID
	events []DomainEventI
}

func (r *AggregateRoot) AddEvent(event DomainEventI) {
	r.events = append(r.events, event)
}
func (r *AggregateRoot) PullEvents() []DomainEventI {
	events := r.events
	r.events = []DomainEventI{}
	return events
}
