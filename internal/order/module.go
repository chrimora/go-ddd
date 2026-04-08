package order

import (
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/order/application/commands"
	"goddd/internal/order/application/queries"
	"goddd/internal/order/domain"
	"goddd/internal/order/infrastructure/sql"
	"goddd/internal/order/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"order_core",
	fx.Provide(
		sql.NewWriteOrderSql,
		sql.NewReadOrderSql,
		fx.Annotate(domain.NewOrderRepository, fx.As(new(domain.OrderRepositoryI))),
		commands.NewCreateOrderCommand,
		commands.NewAddOrderItemCommand,
		queries.NewGetOrderQuery,
		queries.NewGetOrdersByUserQuery,
	),
)

var APIModule = fx.Module(
	"order_api",
	CoreModule,
	fx.Provide(
		commonrest.AsRouteCollection(rest.NewOrderRoutes),
	),
)
