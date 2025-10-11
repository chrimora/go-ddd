package commondomain

import (
	"time"

	"github.com/google/uuid"
)

type AggregateRootI interface {
	PullEvents() []DomainEventI
}
type AggregateRoot struct {
	ID        uuid.UUID
	version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	events    []DomainEventI
}

func NewAggregateRoot() AggregateRoot {
	t := time.Now().UTC()
	return AggregateRoot{ID: uuid.New(), version: 1, CreatedAt: t, UpdatedAt: t}
}

func NewAggregateRootFromFields(id uuid.UUID, version int, createdAt, updatedAt time.Time) AggregateRoot {
	return AggregateRoot{ID: id, version: version, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

func (r *AggregateRoot) GetVersion() int {
	return r.version
}
func (r *AggregateRoot) AddEvent(event DomainEventI) {
	r.events = append(r.events, event)
}
func (r *AggregateRoot) PullEvents() []DomainEventI {
	events := r.events
	r.events = []DomainEventI{}
	return events
}
