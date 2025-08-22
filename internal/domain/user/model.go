package user

import (
	"gotemplate/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type User struct {
	common.AggregateRoot
	ID        uuid.UUID
	UpdatedAt time.Time
	Name      string
}

func NewUser(name string) *User {
	user := &User{
		ID:   uuid.New(),
		Name: name,
	}
	user.AddEvent(UserCreatedEvent{ID: user.ID})
	return user
}

func (u *User) Update(name string) {
	u.Name = name
}
