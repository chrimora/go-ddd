//go:build unit

package domain_test

import (
	"goddd/internal/order/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	order := domain.NewOrder(uuid.New())

	assert.Equal(t, domain.Pending, order.Status())
	assert.Empty(t, order.Items())

	events := order.PullEvents()
	assert.Len(t, events, 1)
	event, ok := events[0].(domain.OrderCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, order.ID(), event.GetAggregateId())
}

func TestAddItem(t *testing.T) {
	t.Run("adds item to pending order", func(t *testing.T) {
		order := domain.NewOrder(uuid.New())

		err := order.AddItem("Widget", 2, 500)

		assert.NoError(t, err)
		assert.Len(t, order.Items(), 1)
		item := order.Items()[0]
		assert.Equal(t, "Widget", item.Name())
		assert.Equal(t, 2, item.Quantity())
		assert.Equal(t, int64(500), item.UnitPrice())
	})

	t.Run("returns error for duplicate item name", func(t *testing.T) {
		order := domain.NewOrder(uuid.New())
		_ = order.AddItem("Widget", 2, 500)

		err := order.AddItem("Widget", 1, 500)

		assert.ErrorIs(t, err, domain.ErrDuplicateItem)
		assert.Len(t, order.Items(), 1)
	})

	t.Run("returns error when order is confirmed", func(t *testing.T) {
		order := domain.NewOrder(uuid.New())
		order.Confirm()

		err := order.AddItem("Widget", 1, 500)

		assert.ErrorIs(t, err, domain.ErrOrderNotPending)
	})

	t.Run("returns error when order is cancelled", func(t *testing.T) {
		order := domain.NewOrder(uuid.New())
		order.Cancel()

		err := order.AddItem("Widget", 1, 500)

		assert.ErrorIs(t, err, domain.ErrOrderNotPending)
	})
}

func TestConfirm(t *testing.T) {
	order := domain.NewOrder(uuid.New())
	order.PullEvents() // clear the created event

	order.Confirm()

	assert.Equal(t, domain.Confirmed, order.Status())
	events := order.PullEvents()
	assert.Len(t, events, 1)
	event, ok := events[0].(domain.OrderConfirmedEvent)
	assert.True(t, ok)
	assert.Equal(t, order.ID(), event.GetAggregateId())
}
