package user

import (
	"fmt"
	"gotemplate/internal/domain/common"
	"time"
)

type User struct {
	common.AggregateRoot
	Name string
}

func NewUser(name string) *User {
	user := &User{
		AggregateRoot: common.NewAggregateRoot(),
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
