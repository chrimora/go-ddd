package test

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/domain"
	"goddd/internal/infrastructure/sql"
	"goddd/test/factories"

	"go.uber.org/fx"
)

var IntegrationTestModule = fx.Module(
	"integration",
	fx.Supply(config.ServiceConfig{Name: "test"}),
	fx.Supply(&config.DBConfig{
		DBHost:     "localhost",
		DBUser:     "goddd",
		DBName:     "goddd",
		DBPassword: "goddd",
	}),
	fx.Provide(
		common.NewLogger,
	),
	domain.Module,
	sql.Module,
	factories.Module,
)
