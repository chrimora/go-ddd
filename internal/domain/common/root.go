package common

import (
	"gotemplate/internal/outbox"

	"github.com/google/uuid"
)

type AggregateRoot struct {
	ID     uuid.UUID
	events []outbox.DomainEventI
}

func (r *AggregateRoot) AddEvent(event outbox.DomainEventI) {
	r.events = append(r.events, event)
}

func (r *AggregateRoot) PullEvents() []outbox.DomainEventI {
	events := r.events
	r.events = []outbox.DomainEventI{}
	return events
}
