package domain

import (
	"fmt"
	"goddd/internal/common/domain"
	"time"
)

type User struct {
	commondomain.AggregateRoot
	Name string
}

func NewUser(name string) *User {
	user := &User{
		AggregateRoot: commondomain.NewAggregateRoot(),
		Name:          name,
	}
	user.AddEvent(NewUserCreatedEvent(user.ID))
	return user
}

func (u *User) String() string {
	return fmt.Sprintf("User[id: %s]", u.ID)
}

func (u *User) Update(name string) {
	u.Name = name
	u.UpdatedAt = time.Now().UTC()
}
