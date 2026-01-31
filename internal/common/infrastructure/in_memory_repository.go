package commoninfrastructure

import (
	"context"
	commondomain "goddd/internal/common/domain"
	"sync"

	"github.com/google/uuid"
)

type Identifiable interface {
	ID() uuid.UUID
}

type InMemoryRepository[T Identifiable] struct {
	mu   sync.RWMutex
	data map[uuid.UUID]T
}

func NewInMemoryRepository[T Identifiable]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		data: make(map[uuid.UUID]T),
	}
}

func (r *InMemoryRepository[T]) Get(
	_ context.Context,
	id uuid.UUID,
) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, ok := r.data[id]
	if !ok {
		var zero T
		return zero, commondomain.ErrNotFound
	}

	return clone(entity), nil
}

func (r *InMemoryRepository[T]) Create(
	_ context.Context,
	entity T,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.ID()
	r.data[id] = clone(entity)
	return nil
}

func (r *InMemoryRepository[T]) Update(
	_ context.Context,
	entity T,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.ID()
	if _, exists := r.data[id]; !exists {
		return commondomain.ErrNotFound
	}

	r.data[id] = clone(entity)
	return nil
}

func (r *InMemoryRepository[T]) Remove(
	_ context.Context,
	entity T,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.ID()
	if _, exists := r.data[id]; !exists {
		return commondomain.ErrNotFound
	}

	delete(r.data, id)
	return nil
}

func clone[T any](v T) T {
	return v
}
