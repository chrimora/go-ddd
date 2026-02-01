package commondomain

import (
	"time"

	"github.com/google/uuid"
)

type AggregateRoot struct {
	id        uuid.UUID
	version   int
	createdAt time.Time
	updatedAt time.Time
	events    []DomainEventI
}

func NewAggregateRoot() AggregateRoot {
	t := time.Now().UTC()
	return AggregateRoot{id: NewUUID(), version: 0, createdAt: t, updatedAt: t}
}
func RehydrateAggregateRoot(
	id uuid.UUID, version int, createdAt, updatedAt time.Time,
) AggregateRoot {
	return AggregateRoot{id: id, version: version, createdAt: createdAt, updatedAt: updatedAt}
}

// Getters
func (r *AggregateRoot) ID() uuid.UUID        { return r.id }
func (r *AggregateRoot) Version() int         { return r.version }
func (r *AggregateRoot) CreatedAt() time.Time { return r.createdAt }
func (r *AggregateRoot) UpdatedAt() time.Time { return r.updatedAt }

func (r *AggregateRoot) Update() {
	r.updatedAt = time.Now().UTC()
}
func (r *AggregateRoot) AddEvent(event DomainEventI) {
	r.events = append(r.events, event)
}
func (r *AggregateRoot) PullEvents() []DomainEventI {
	events := r.events
	r.events = []DomainEventI{}
	return events
}
