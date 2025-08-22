package user

import (
	"gotemplate/internal/outbox"

	"github.com/google/uuid"
)

const UserCreated outbox.EventType = "userCreated"

type UserCreatedEvent struct {
	ID uuid.UUID
}

func (u UserCreatedEvent) Type() outbox.EventType {
	return UserCreated
}
