package domain

import (
	"fmt"
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "Pending"
	OrderConfirmed OrderStatus = "Confirmed"
	OrderCancelled OrderStatus = "Cancelled"
)

type Order struct {
	commondomain.AggregateRoot
	status OrderStatus
	total  int64
}

func NewOrder(total int64) *Order {
	order := &Order{
		AggregateRoot: commondomain.NewAggregateRoot(),
		status:        OrderPending,
		total:         total,
	}
	order.AddEvent(NewOrderCreatedEvent(order.ID()))
	return order
}

func RehydrateOrder(id uuid.UUID, version int, status OrderStatus, total int64) *Order {
	return &Order{
		AggregateRoot: commondomain.RehydrateAggregateRoot(id, version),
		status:        status,
		total:         total,
	}
}

func (o *Order) Clone() *Order {
	return &Order{
		AggregateRoot: o.AggregateRoot.Clone(),
		status:        o.status,
		total:         o.total,
	}
}

func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Total() int64        { return o.total }
func (o *Order) String() string      { return fmt.Sprintf("Order[id: %s]", o.ID()) }

func (o *Order) Confirm() {
	o.status = OrderConfirmed
}

func (o *Order) Cancel() {
	o.status = OrderCancelled
}
