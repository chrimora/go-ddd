package commontest

import (
	"goddd/internal/common"
	"goddd/internal/common/infrastructure"
	"goddd/internal/config"

	"go.uber.org/fx"
)

var UnitTestModule = fx.Module(
	"test_unit",
	fx.Supply(config.ServiceConfig{Name: "test"}),
	fx.Provide(
		commoninfrastructure.NewLogger,
	),
)

var IntegrationTestModule = fx.Module(
	"test_integration",
	fx.Supply(config.ServiceConfig{Name: "test"}),
	fx.Supply(&config.DBConfig{
		DBHost:     "localhost",
		DBUser:     "goddd",
		DBName:     "goddd",
		DBPassword: "goddd",
	}),
	common.CoreModule,
)
