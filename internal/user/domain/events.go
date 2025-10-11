package domain

import (
	"goddd/internal/common/domain"

	"github.com/google/uuid"
)

type userEventRoot struct{ commondomain.DomainEventRoot }

func (e userEventRoot) GetAggregateType() string {
	return "User"
}
func newUserEventRoot(aggregateId uuid.UUID) userEventRoot {
	return userEventRoot{
		DomainEventRoot: commondomain.DomainEventRoot{AggregateId: aggregateId},
	}
}

const UserCreated commondomain.EventType = "userCreated"

type UserCreatedEvent struct{ userEventRoot }

func (e UserCreatedEvent) GetEventType() commondomain.EventType {
	return UserCreated
}
func NewUserCreatedEvent(aggregateId uuid.UUID) UserCreatedEvent {
	return UserCreatedEvent{
		userEventRoot: newUserEventRoot(aggregateId),
	}
}
