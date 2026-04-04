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

type CreateOrderInput struct {
	UserId uuid.UUID
}

type CreateOrderCommand commonapplication.CommandI[CreateOrderInput]

func NewCreateOrderCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	orderRepo domain.OrderRepositoryI,
) CreateOrderCommand {
	return commonapplication.NewCommand(log, &createOrder{
		txManager: txManager,
		orderRepo: orderRepo,
	})
}

type createOrder struct {
	txManager commondomain.TxManager
	orderRepo domain.OrderRepositoryI
}

func (c *createOrder) Handle(
	ctx context.Context, log *slog.Logger, input CreateOrderInput,
) (uuid.UUID, error) {
	order := domain.NewOrder(input.UserId)
	log.InfoContext(ctx, "Creating", "order", order)

	err := c.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return c.orderRepo.Create(ctx, tx, order)
	})
	return order.ID(), err
}
