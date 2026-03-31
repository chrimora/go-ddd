package domain

import (
	"context"
	"database/sql"
	"errors"
	"goddd/internal/outbox"
	ordersql "goddd/internal/order/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderRepositoryI interface {
	Get(context.Context, uuid.UUID) (*Order, error)
	Create(context.Context, pgx.Tx, *Order) error
	Update(context.Context, pgx.Tx, *Order) error
	Remove(context.Context, pgx.Tx, *Order) error
}

type OrderRepository OrderRepositoryI

type orderRepository struct {
	log        *slog.Logger
	orderSql   ordersql.WriteOrderSql
	outboxRepo outbox.OutboxRepositoryI
}

func NewOrderRepository(
	log *slog.Logger,
	orderSql ordersql.WriteOrderSql,
	outboxRepo outbox.OutboxRepositoryI,
) OrderRepository {
	return &orderRepository{
		log:        log,
		orderSql:   orderSql,
		outboxRepo: outboxRepo,
	}
}

func (r *orderRepository) Get(ctx context.Context, id uuid.UUID) (*Order, error) {
	row, err := r.orderSql.GetOrder(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound(id)
		default:
			return nil, err
		}
	}
	return RehydrateOrder(row.ID, int(row.Version), OrderStatus(row.Status), row.Total), nil
}

func (r *orderRepository) Create(ctx context.Context, tx pgx.Tx, order *Order) error {
	_, err := r.orderSql.WithTx(tx).CreateOrder(
		ctx,
		ordersql.CreateOrderParams{
			ID:      order.ID(),
			Version: int32(order.Version()),
			Status:  string(order.Status()),
			Total:   order.Total(),
		},
	)
	if err != nil {
		return err
	}
	return r.outboxRepo.CreateMany(ctx, tx, order.PullEvents()...)
}

func (r *orderRepository) Update(ctx context.Context, tx pgx.Tx, order *Order) error {
	_, err := r.orderSql.WithTx(tx).UpdateOrder(
		ctx,
		ordersql.UpdateOrderParams{
			ID:      order.ID(),
			Version: int32(order.Version()),
			Status:  string(order.Status()),
			Total:   order.Total(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRaceCondition(order.ID())
		default:
			return err
		}
	}
	return r.outboxRepo.CreateMany(ctx, tx, order.PullEvents()...)
}

func (r *orderRepository) Remove(ctx context.Context, tx pgx.Tx, order *Order) error {
	return r.orderSql.WithTx(tx).RemoveOrder(ctx, order.ID())
}
