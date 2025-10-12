package commontest

import "testing"

type Mock[T any] interface {
	Mock(t *testing.T, overrides ...map[string]any) *T
}
