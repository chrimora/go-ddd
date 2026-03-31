//go:build unit

package domain_test

import (
	"goddd/internal/user/domain"
	"testing"
	"time"

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

func TestUserUpdate(t *testing.T) {
	user := domain.NewUser("Bob")
	oldUpdatedAt := user.UpdatedAt()

	// Ensure time difference
	time.Sleep(1 * time.Millisecond)

	newName := "Robert"
	user.ChangeName(newName)

	assert.Equal(t, newName, user.Name())
	assert.Greater(t, user.UpdatedAt(), oldUpdatedAt)
}
