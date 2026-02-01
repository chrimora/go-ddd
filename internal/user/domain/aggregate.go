package domain

import (
	"fmt"
	"goddd/internal/common/domain"
	"time"

	"github.com/google/uuid"
)

type User struct {
	commondomain.AggregateRoot
	name string
}

func NewUser(name string) *User {
	user := &User{
		AggregateRoot: commondomain.NewAggregateRoot(),
		name:          name,
	}
	user.AddEvent(NewUserCreatedEvent(user.ID()))
	return user
}
func RehydrateUser(
	id uuid.UUID, version int, createdAt, updatedAt time.Time, name string,
) *User {
	root := commondomain.RehydrateAggregateRoot(id, version, createdAt, updatedAt)
	return &User{
		AggregateRoot: root,
		name:          name,
	}
}

// Getters
func (u *User) Name() string   { return u.name }
func (u *User) String() string { return fmt.Sprintf("User[id: %s]", u.ID()) }

func (u *User) ChangeName(name string) {
	u.name = name
	u.AggregateRoot.Update()
}
