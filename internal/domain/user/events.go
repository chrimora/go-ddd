package user

import (
	"gotemplate/internal/domain/common"

	"github.com/google/uuid"
)

const UserCreated common.EventType = "userCreated"

type UserCreatedEvent struct {
	common.DomainEventRoot
}

func NewUserCreatedEvent(aggregateId uuid.UUID) UserCreatedEvent {
	return UserCreatedEvent{
		DomainEventRoot: common.DomainEventRoot{
			EventType: UserCreated, AggregateId: aggregateId,
		},
	}
}
