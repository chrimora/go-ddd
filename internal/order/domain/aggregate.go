package domain

import (
	"fmt"
	commondomain "goddd/internal/common/domain"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	Pending   OrderStatus = "Pending"
	Confirmed OrderStatus = "Confirmed"
	Cancelled OrderStatus = "Cancelled"
)

type OrderItem struct {
	id        uuid.UUID
	name      string
	quantity  int
	unitPrice int64 // pence
}

func NewOrderItem(name string, quantity int, unitPrice int64) OrderItem {
	return OrderItem{
		id:        commondomain.NewUUID(),
		name:      name,
		quantity:  quantity,
		unitPrice: unitPrice,
	}
}

func RehydrateOrderItem(id uuid.UUID, name string, quantity int, unitPrice int64) OrderItem {
	return OrderItem{id: id, name: name, quantity: quantity, unitPrice: unitPrice}
}

func (i OrderItem) ID() uuid.UUID    { return i.id }
func (i OrderItem) Name() string     { return i.name }
func (i OrderItem) Quantity() int    { return i.quantity }
func (i OrderItem) UnitPrice() int64 { return i.unitPrice }

type Order struct {
	commondomain.AggregateRoot
	userId uuid.UUID
	status OrderStatus
	items  []OrderItem
}

func NewOrder(userId uuid.UUID) *Order {
	order := &Order{
		AggregateRoot: commondomain.NewAggregateRoot(),
		userId:        userId,
		status:        Pending,
		items:         []OrderItem{},
	}
	order.AddEvent(NewOrderCreatedEvent(order.ID()))
	return order
}

func RehydrateOrder(id uuid.UUID, version int, userId uuid.UUID, status OrderStatus, items []OrderItem) *Order {
	return &Order{
		AggregateRoot: commondomain.RehydrateAggregateRoot(id, version),
		userId:        userId,
		status:        status,
		items:         items,
	}
}

func (o *Order) UserID() uuid.UUID   { return o.userId }
func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Items() []OrderItem  { return o.items }
func (o *Order) String() string      { return fmt.Sprintf("Order[id: %s]", o.ID()) }

func (o *Order) AddItem(name string, quantity int, unitPrice int64) error {
	if o.status != Pending {
		return ErrOrderNotPending
	}
	o.items = append(o.items, NewOrderItem(name, quantity, unitPrice))
	return nil
}

func (o *Order) Confirm() {
	o.status = Confirmed
	o.AddEvent(NewOrderConfirmedEvent(o.ID()))
}

func (o *Order) Cancel() {
	o.status = Cancelled
}

func (o *Order) Clone() *Order {
	items := make([]OrderItem, len(o.items))
	copy(items, o.items)
	return &Order{
		AggregateRoot: o.AggregateRoot.Clone(),
		userId:        o.userId,
		status:        o.status,
		items:         items,
	}
}
