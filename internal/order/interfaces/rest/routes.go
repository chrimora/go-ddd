package rest

import (
	"context"
	"errors"
	commondomain "goddd/internal/common/domain"
	commoninfrastructure "goddd/internal/common/infrastructure"
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/order/application/commands"
	"goddd/internal/order/application/queries"
	"goddd/internal/order/domain"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type (
	OrderRoutes commonrest.RouteCollection
	orderRoutes struct {
		log             *slog.Logger
		createOrder     commands.CreateOrderCommand
		addOrderItem    commands.AddOrderItemCommand
		getOrder        queries.GetOrderQuery
		getOrdersByUser queries.GetOrdersByUserQuery
	}
)

func NewOrderRoutes(
	log *slog.Logger,
	createOrder commands.CreateOrderCommand,
	addOrderItem commands.AddOrderItemCommand,
	getOrder queries.GetOrderQuery,
	getOrdersByUser queries.GetOrdersByUserQuery,
) OrderRoutes {
	return &orderRoutes{
		log:             log,
		createOrder:     createOrder,
		addOrderItem:    addOrderItem,
		getOrder:        getOrder,
		getOrdersByUser: getOrdersByUser,
	}
}

func (o *orderRoutes) Register(api huma.API) {
	huma.Post(api, "/orders", o.create)
	huma.Get(api, "/orders/{id}", o.get)
	huma.Get(api, "/orders", o.getByUser)
	huma.Post(api, "/orders/{id}/items", o.addItem)
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

type OrderSummaryPayload struct {
	commonrest.IdPayload
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type ListOrdersQuery struct {
	commonrest.PaginationQuery
}

func (o *orderRoutes) getByUser(
	ctx context.Context, req *ListOrdersQuery,
) (*commonrest.Response[commonrest.Page[OrderSummaryPayload]], error) {
	rc := commoninfrastructure.MustGetRequestCtx(ctx)

	var after *uuid.UUID
	if req.AfterCursor != uuid.Nil {
		after = &req.AfterCursor
	}

	out, err := o.getOrdersByUser.Handle(ctx, queries.GetOrdersByUserInput{
		UserId: rc.UserId,
		Limit:  req.Limit,
		After:  after,
	})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(o.log, ctx, err)
	}

	items := make([]OrderSummaryPayload, len(out.Orders))
	for i, s := range out.Orders {
		items[i] = OrderSummaryPayload{Status: s.Status, Total: s.Total}
		items[i].ID = s.ID
	}
	return commonrest.BuildResponse(commonrest.Page[OrderSummaryPayload]{
		Items:      items,
		NextCursor: out.Next,
	}), nil
}

type AddItemPayload struct {
	Name     string `json:"name" minLength:"1"`
	Quantity int    `json:"quantity" minimum:"1"`
}

func (o *orderRoutes) addItem(
	ctx context.Context, req *commonrest.UpdateRequest[AddItemPayload],
) (*commonrest.EmptyResponse, error) {
	_, err := o.addOrderItem.Handle(ctx, commands.AddOrderItemInput{
		OrderId:  req.ID,
		Name:     req.Body.Name,
		Quantity: req.Body.Quantity,
	})
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(o.log, ctx, err)
		case errors.Is(err, domain.ErrOrderNotPending):
			return nil, huma.Error409Conflict(err.Error())
		case errors.Is(err, domain.ErrDuplicateItem):
			return nil, huma.Error409Conflict(err.Error())
		default:
			return nil, commonrest.UnexpectedErrorResponse(o.log, ctx, err)
		}
	}
	return &commonrest.EmptyResponse{}, nil
}
