package domain

import (
	"goddd/internal/common/domain"

	"github.com/google/uuid"
)

type userEventRoot struct{ domain.DomainEventRoot }

func (e userEventRoot) GetAggregateType() string {
	return "User"
}
func newUserEventRoot(aggregateId uuid.UUID) userEventRoot {
	return userEventRoot{
		DomainEventRoot: domain.DomainEventRoot{AggregateId: aggregateId},
	}
}

const UserCreated domain.EventType = "userCreated"

type UserCreatedEvent struct{ userEventRoot }

func (e UserCreatedEvent) GetEventType() domain.EventType {
	return UserCreated
}
func NewUserCreatedEvent(aggregateId uuid.UUID) UserCreatedEvent {
	return UserCreatedEvent{
		userEventRoot: newUserEventRoot(aggregateId),
	}
}
