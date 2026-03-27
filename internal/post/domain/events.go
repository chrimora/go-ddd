package domain

import (
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

type postEventRoot struct{ commondomain.DomainEventRoot }

func (e postEventRoot) GetAggregateType() string {
	return "Post"
}
func newPostEventRoot(aggregateId uuid.UUID) postEventRoot {
	return postEventRoot{
		DomainEventRoot: commondomain.DomainEventRoot{AggregateId: aggregateId},
	}
}

const PostCreated commondomain.EventType = "postCreated"

type PostCreatedEvent struct{ postEventRoot }

func (e PostCreatedEvent) GetEventType() commondomain.EventType { return PostCreated }
func NewPostCreatedEvent(aggregateId uuid.UUID) PostCreatedEvent {
	return PostCreatedEvent{postEventRoot: newPostEventRoot(aggregateId)}
}

const PostUpdated commondomain.EventType = "postUpdated"

type PostUpdatedEvent struct{ postEventRoot }

func (e PostUpdatedEvent) GetEventType() commondomain.EventType { return PostUpdated }
func NewPostUpdatedEvent(aggregateId uuid.UUID) PostUpdatedEvent {
	return PostUpdatedEvent{postEventRoot: newPostEventRoot(aggregateId)}
}
