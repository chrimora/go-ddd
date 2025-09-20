package user

import (
	"fmt"
	"gotemplate/internal/domain/common"

	"github.com/google/uuid"
)

func ErrNotFound(id uuid.UUID) error {
	return fmt.Errorf("%w: user %s", common.ErrNotFound, id)
}
func ErrRaceCondition(id uuid.UUID) error {
	return fmt.Errorf("%w: user %s", common.ErrRaceCondition, id)
}
