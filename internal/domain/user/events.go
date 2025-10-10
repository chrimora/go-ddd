package user

import (
	"goddd/internal/domain/common"

	"github.com/google/uuid"
)

type userEventRoot struct{ common.DomainEventRoot }

func (e userEventRoot) GetAggregateType() string {
	return "User"
}
func newUserEventRoot(aggregateId uuid.UUID) userEventRoot {
	return userEventRoot{
		DomainEventRoot: common.DomainEventRoot{AggregateId: aggregateId},
	}
}

const UserCreated common.EventType = "userCreated"

type UserCreatedEvent struct{ userEventRoot }

func (e UserCreatedEvent) GetEventType() common.EventType {
	return UserCreated
}
func NewUserCreatedEvent(aggregateId uuid.UUID) UserCreatedEvent {
	return UserCreatedEvent{
		userEventRoot: newUserEventRoot(aggregateId),
	}
}
