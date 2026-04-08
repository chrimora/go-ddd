package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/order/domain"
	"log/slog"

	"github.com/google/uuid"
)

type GetOrderInput struct {
	Id uuid.UUID
}

type GetOrderQuery commonapplication.QueryI[GetOrderInput, *domain.Order]

func NewGetOrderQuery(log *slog.Logger, orderRepo domain.OrderRepositoryI) GetOrderQuery {
	return commonapplication.NewQuery(log, &getOrder{orderRepo: orderRepo})
}

type getOrder struct {
	orderRepo domain.OrderRepositoryI
}

func (q *getOrder) Handle(ctx context.Context, log *slog.Logger, input GetOrderInput) (*domain.Order, error) {
	order, err := q.orderRepo.Get(ctx, input.Id)
	log.InfoContext(ctx, "Got", "order", order)
	return order, err
}
