package domain

import (
	"errors"
	"fmt"
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

var (
	ErrOrderNotPending  = errors.New("order is not pending")
	ErrDuplicateItem    = errors.New("item with that name already exists")
)

func ErrNotFound(id uuid.UUID) error {
	return fmt.Errorf("%w: order %s", commondomain.ErrNotFound, id)
}
func ErrRaceCondition(id uuid.UUID) error {
	return fmt.Errorf("%w: order %s", commondomain.ErrRaceCondition, id)
}
