package test

import (
	commoninfrastructure "goddd/internal/common/infrastructure"
	"goddd/internal/common/test"
	"goddd/internal/order"
	"goddd/internal/order/application/commands"
	"goddd/internal/order/application/queries"
	"goddd/internal/order/domain"

	"go.uber.org/fx"
)

var UnitTestModule = fx.Module(
	"order_unit_test",
	commontest.UnitTestModule,
	fx.Provide(
		fx.Annotate(
			commoninfrastructure.NewInMemoryRepository[*domain.Order],
			fx.As(new(domain.OrderRepositoryI)),
		),
		NewOrderFactory,
		commands.NewCreateOrderCommand,
		commands.NewAddOrderItemCommand,
		queries.NewGetOrderQuery,
	),
)

var IntegrationTestModule = fx.Module(
	"order_integration_test",
	commontest.IntegrationTestModule,
	order.CoreModule,
	fx.Provide(
		NewOrderFactory,
	),
)
