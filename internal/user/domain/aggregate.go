package domain

import (
	"fmt"
	"goddd/internal/common/domain"

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

func RehydrateUser(id uuid.UUID, version int, name string) *User {
	return &User{
		AggregateRoot: commondomain.RehydrateAggregateRoot(id, version),
		name:          name,
	}
}

// Getters
func (u *User) Name() string   { return u.name }
func (u *User) String() string { return fmt.Sprintf("User[id: %s]", u.ID()) }

func (u *User) ChangeName(name string) {
	u.name = name
}

func (u *User) Clone() *User {
	return &User{
		AggregateRoot: u.AggregateRoot.Clone(),
		name:          u.name,
	}
}
