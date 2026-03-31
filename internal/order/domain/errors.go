package domain

import (
	"fmt"
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

func ErrNotFound(id uuid.UUID) error {
	return fmt.Errorf("%w: order %s", commondomain.ErrNotFound, id)
}
func ErrRaceCondition(id uuid.UUID) error {
	return fmt.Errorf("%w: order %s", commondomain.ErrRaceCondition, id)
}
