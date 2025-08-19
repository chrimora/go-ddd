package user

import (
	userdb "gotemplate/internal/user/db"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	// Events    []events.Event
	Name string
}

func fromDB(user userdb.User) *User {
	return &User{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}
}
