package rest

import (
	"context"
	"errors"
	commoninfrastructure "goddd/internal/common/infrastructure"
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/order/application/commands"
	"goddd/internal/order/application/queries"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/order/domain"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type (
	OrderRoutes commonrest.RouteCollection
	orderRoutes struct {
		log         *slog.Logger
		createOrder commands.CreateOrderCommand
		getOrder    queries.GetOrderQuery
	}
)

func NewOrderRoutes(
	log *slog.Logger,
	createOrder commands.CreateOrderCommand,
	getOrder queries.GetOrderQuery,
) OrderRoutes {
	return &orderRoutes{
		log:         log,
		createOrder: createOrder,
		getOrder:    getOrder,
	}
}

func (o *orderRoutes) Register(api huma.API) {
	huma.Post(api, "/orders", o.create)
	huma.Get(api, "/orders/{id}", o.get)
}

func (o *orderRoutes) create(
	ctx context.Context, req *commonrest.CreateRequest[struct{}],
) (*commonrest.Response[commonrest.IdPayload], error) {
	rc := commoninfrastructure.MustGetRequestCtx(ctx)

	id, err := o.createOrder.Handle(ctx, commands.CreateOrderInput{UserId: rc.UserId})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(o.log, ctx, err)
	}
	return commonrest.BuildResponse(commonrest.IdPayload{ID: id}), nil
}

type OrderItemPayload struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
}

type OrderPayload struct {
	commonrest.IdPayload
	UserID uuid.UUID          `json:"user_id"`
	Status domain.OrderStatus `json:"status"`
	Items  []OrderItemPayload `json:"items"`
}

func (o *orderRoutes) get(
	ctx context.Context, req *commonrest.IdParam,
) (*commonrest.Response[OrderPayload], error) {
	order, err := o.getOrder.Handle(ctx, queries.GetOrderInput{Id: req.ID})
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(o.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(o.log, ctx, err)
		}
	}

	items := make([]OrderItemPayload, len(order.Items()))
	for i, item := range order.Items() {
		items[i] = OrderItemPayload{
			ID:        item.ID(),
			Name:      item.Name(),
			Quantity:  item.Quantity(),
			UnitPrice: item.UnitPrice(),
		}
	}

	res := OrderPayload{
		UserID: order.UserID(),
		Status: order.Status(),
		Items:  items,
	}
	res.ID = order.ID()
	return commonrest.BuildResponse(res), nil
}
