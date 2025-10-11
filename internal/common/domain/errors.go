package commondomain

import (
	"errors"
	"fmt"
)

var (
	errShouldNotHappen = errors.New("should not happen")

	ErrNotFound      = errors.New("not found")
	ErrRaceCondition = errors.New("race condition")
)

func ErrShouldNotHappen(msg string) error {
	return fmt.Errorf("%w: %s", errShouldNotHappen, msg)
}
