package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	ordersql "goddd/internal/order/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
)

type OrderSummary struct {
	ID     uuid.UUID
	Status string
	Total  int64
}

type GetOrdersByUserInput struct {
	UserId uuid.UUID
	Limit  int
	After  *uuid.UUID
}

type GetOrdersByUserOutput struct {
	Orders []OrderSummary
	Next   *uuid.UUID
}

type GetOrdersByUserQuery commonapplication.QueryI[GetOrdersByUserInput, GetOrdersByUserOutput]

func NewGetOrdersByUserQuery(log *slog.Logger, readSql ordersql.ReadOrderSql) GetOrdersByUserQuery {
	return commonapplication.NewQuery(log, &getOrdersByUser{readSql: readSql})
}

type getOrdersByUser struct {
	readSql ordersql.ReadOrderSql
}

func (q *getOrdersByUser) Handle(ctx context.Context, log *slog.Logger, input GetOrdersByUserInput) (GetOrdersByUserOutput, error) {
	rows, err := q.readSql.GetOrderSummariesByUserId(ctx, ordersql.GetOrderSummariesByUserIdParams{
		UserID:       input.UserId,
		After:        input.After,
		LimitPlusOne: int32(input.Limit + 1),
	})
	if err != nil {
		return GetOrdersByUserOutput{}, err
	}

	var next *uuid.UUID
	if len(rows) > input.Limit {
		rows = rows[:input.Limit]
		next = &rows[input.Limit-1].ID
	}

	summaries := make([]OrderSummary, len(rows))
	for i, r := range rows {
		summaries[i] = OrderSummary{ID: r.ID, Status: r.Status, Total: r.Total}
	}
	return GetOrdersByUserOutput{Orders: summaries, Next: next}, nil
}
