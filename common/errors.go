package common

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrRaceCondition = errors.New("race condition")
)
