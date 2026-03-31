package commondomain

import (
	"github.com/google/uuid"
)

type AggregateRoot struct {
	id      uuid.UUID
	version int
	events  []DomainEventI
}

func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{id: NewUUID(), version: 0}
}

func RehydrateAggregateRoot(id uuid.UUID, version int) AggregateRoot {
	return AggregateRoot{id: id, version: version}
}

// Getters
func (r *AggregateRoot) ID() uuid.UUID { return r.id }
func (r *AggregateRoot) Version() int  { return r.version }

func (r *AggregateRoot) AddEvent(event DomainEventI) {
	r.events = append(r.events, event)
}
func (r *AggregateRoot) PullEvents() []DomainEventI {
	events := r.events
	r.events = []DomainEventI{}
	return events
}

func (r *AggregateRoot) Clone() AggregateRoot {
	events := make([]DomainEventI, len(r.events))
	copy(events, r.events)
	return AggregateRoot{
		id:      r.id,
		version: r.version,
		events:  events,
	}
}
