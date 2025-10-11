package domain

import (
	"fmt"
	"goddd/internal/common/domain"

	"github.com/google/uuid"
)

func ErrNotFound(id uuid.UUID) error {
	return fmt.Errorf("%w: user %s", commondomain.ErrNotFound, id)
}
func ErrRaceCondition(id uuid.UUID) error {
	return fmt.Errorf("%w: user %s", commondomain.ErrRaceCondition, id)
}
