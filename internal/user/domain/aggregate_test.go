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

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, name, user.Name())
}

func TestUserUpdate(t *testing.T) {
	user := domain.NewUser("Bob")
	oldUpdatedAt := user.UpdatedAt()

	// Ensure time difference
	time.Sleep(3 * time.Millisecond)

	newName := "Robert"
	user.Update(newName)

	assert.Equal(t, newName, user.Name())
	assert.Greater(t, user.UpdatedAt(), oldUpdatedAt)
}
