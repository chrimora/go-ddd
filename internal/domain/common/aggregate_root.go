package common

import (
	"github.com/google/uuid"
)

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
