package test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/common/test"
	"goddd/internal/order/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrderFactory commontest.Mock[domain.Order]

type orderFactory struct {
	repo      domain.OrderRepositoryI
	txManager commondomain.TxManager
}

func NewOrderFactory(repo domain.OrderRepositoryI, txManager commondomain.TxManager) OrderFactory {
	return &orderFactory{repo: repo, txManager: txManager}
}

func (f *orderFactory) Mock(t *testing.T, ctx context.Context, overrides ...map[string]any) *domain.Order {
	fields := &struct {
		ID     uuid.UUID
		Version int
		UserId uuid.UUID
		Status domain.OrderStatus
	}{
		ID:     commondomain.NewUUID(),
		UserId: commondomain.NewUUID(),
		Status: domain.Pending,
	}
	commontest.Merge(fields, overrides)

	order := domain.RehydrateOrder(fields.ID, fields.Version, fields.UserId, fields.Status, []domain.OrderItem{})
	err := f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return f.repo.Create(ctx, tx, order)
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		err := f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
			return f.repo.Remove(ctx, tx, order)
		})
		if err != nil {
			panic(err)
		}
	})
	return order
}
