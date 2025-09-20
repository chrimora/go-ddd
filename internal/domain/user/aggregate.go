package user

import (
	"fmt"
	"gotemplate/internal/domain/common"
	"time"

	"github.com/google/uuid"
)

type User struct {
	common.AggregateRoot
	UpdatedAt time.Time
	Name      string
}

func NewUser(name string) *User {
	user := &User{
		AggregateRoot: common.AggregateRoot{ID: uuid.New()},
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
}
