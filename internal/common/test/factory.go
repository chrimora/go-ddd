package commontest

import (
	"context"
	"testing"
)

type Mock[T any] interface {
	Mock(t *testing.T, ctx context.Context, overrides ...map[string]any) *T
}
