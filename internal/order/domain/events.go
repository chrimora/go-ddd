package domain

import (
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

type orderEventRoot struct{ commondomain.DomainEventRoot }

func (e orderEventRoot) GetAggregateType() string {
	return "Order"
}
func newOrderEventRoot(aggregateId uuid.UUID) orderEventRoot {
	return orderEventRoot{
		DomainEventRoot: commondomain.DomainEventRoot{AggregateId: aggregateId},
	}
}

const OrderCreated commondomain.EventType = "orderCreated"

type OrderCreatedEvent struct{ orderEventRoot }

func (e OrderCreatedEvent) GetEventType() commondomain.EventType {
	return OrderCreated
}
func NewOrderCreatedEvent(aggregateId uuid.UUID) OrderCreatedEvent {
	return OrderCreatedEvent{
		orderEventRoot: newOrderEventRoot(aggregateId),
	}
}

const OrderConfirmed commondomain.EventType = "orderConfirmed"

type OrderConfirmedEvent struct{ orderEventRoot }

func (e OrderConfirmedEvent) GetEventType() commondomain.EventType {
	return OrderConfirmed
}
func NewOrderConfirmedEvent(aggregateId uuid.UUID) OrderConfirmedEvent {
	return OrderConfirmedEvent{
		orderEventRoot: newOrderEventRoot(aggregateId),
	}
}
