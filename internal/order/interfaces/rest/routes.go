package rest

import (
	"context"
	commoninfrastructure "goddd/internal/common/infrastructure"
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/order/application/commands"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type (
	OrderRoutes commonrest.RouteCollection
	orderRoutes struct {
		log         *slog.Logger
		createOrder commands.CreateOrderCommand
	}
)

func NewOrderRoutes(
	log *slog.Logger,
	createOrder commands.CreateOrderCommand,
) OrderRoutes {
	return &orderRoutes{
		log:         log,
		createOrder: createOrder,
	}
}

func (o *orderRoutes) Register(api huma.API) {
	huma.Post(api, "/orders", o.create)
}

func (o *orderRoutes) create(
	ctx context.Context, req *commonrest.CreateRequest[struct{}],
) (*commonrest.Response[commonrest.IdPayload], error) {
	rc := commoninfrastructure.MustGetRequestCtx(ctx)
	userId, err := uuid.Parse(rc.UserId)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid user id in request context")
	}

	id, err := o.createOrder.Handle(ctx, commands.CreateOrderInput{UserId: userId})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(o.log, ctx, err)
	}
	return commonrest.BuildResponse(commonrest.IdPayload{ID: id}), nil
}
