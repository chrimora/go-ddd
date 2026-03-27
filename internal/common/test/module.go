package commontest

import (
	"goddd/internal/common"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	"goddd/internal/config"
	"goddd/internal/outbox"

	"go.uber.org/fx"
)

var UnitTestModule = fx.Module(
	"test_unit",
	fx.Supply(&config.ServiceConfig{Name: "test", Env: config.TestEnvEnum}),
	fx.Provide(
		commoninfrastructure.NewLogger,
		fx.Annotate(commondomain.NewMockTxManager, fx.As(new(commondomain.TxManager))),
	),
)

var IntegrationTestModule = fx.Module(
	"test_integration",
	fx.Supply(&config.ServiceConfig{Name: "test", Env: config.TestEnvEnum}),
	fx.Provide(config.NewConfig[config.DBConfig]),
	common.CoreModule,
	outbox.CoreModule,
)
