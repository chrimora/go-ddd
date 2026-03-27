package commoninfrastructure

import (
	"context"
	commondomain "goddd/internal/common/domain"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Identifiable[T any] interface {
	ID() uuid.UUID
	Clone() T
}

type InMemoryRepository[T Identifiable[T]] struct {
	mu   sync.RWMutex
	data map[uuid.UUID]T
}

func NewInMemoryRepository[T Identifiable[T]]() *InMemoryRepository[T] {
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

	return entity.Clone(), nil
}

func (r *InMemoryRepository[T]) Create(
	_ context.Context,
	_ pgx.Tx,
	entity T,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.ID()
	r.data[id] = entity.Clone()
	return nil
}

func (r *InMemoryRepository[T]) Update(
	_ context.Context,
	_ pgx.Tx,
	entity T,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := entity.ID()
	if _, exists := r.data[id]; !exists {
		return commondomain.ErrNotFound
	}

	r.data[id] = entity.Clone()
	return nil
}

func (r *InMemoryRepository[T]) Remove(
	_ context.Context,
	_ pgx.Tx,
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
