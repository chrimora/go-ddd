package commands

import (
	"context"
	commonapplication "goddd/internal/common/application"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/order/domain"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AddOrderItemInput struct {
	OrderId  uuid.UUID
	Name     string
	Quantity int
}

type AddOrderItemCommand commonapplication.CommandI[AddOrderItemInput]

func NewAddOrderItemCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	orderRepo domain.OrderRepositoryI,
) AddOrderItemCommand {
	return commonapplication.NewCommand(log, &addOrderItem{
		txManager: txManager,
		orderRepo: orderRepo,
	})
}

type addOrderItem struct {
	txManager commondomain.TxManager
	orderRepo domain.OrderRepositoryI
}

func (c *addOrderItem) Handle(
	ctx context.Context, log *slog.Logger, input AddOrderItemInput,
) (uuid.UUID, error) {
	// In real life the unit price would be looked up
	unitPrice := 1000

	order, err := c.orderRepo.Get(ctx, input.OrderId)
	if err != nil {
		return uuid.Nil, err
	}
	log.InfoContext(ctx, "Adding item", "order", order, "name", input.Name)
	if err := order.AddItem(input.Name, input.Quantity, int64(unitPrice)); err != nil {
		return uuid.Nil, err
	}
	err = c.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return c.orderRepo.Update(ctx, tx, order)
	})
	return input.OrderId, err
}
