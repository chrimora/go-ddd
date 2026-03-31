//go:build unit

package domain_test

import (
	"goddd/internal/user/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	name := "Alice"
	user := domain.NewUser(name)

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, name, user.Name())

	events := user.PullEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(domain.UserCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), event.GetAggregateId())
}
