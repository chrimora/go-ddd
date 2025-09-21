package test

import (
	"gotemplate/internal/common"
	"gotemplate/internal/config"
	"gotemplate/internal/domain"
	"gotemplate/internal/infrastructure/sql"
	"gotemplate/test/factories"

	"go.uber.org/fx"
)

var IntegrationTestModule = fx.Module(
	"integration",
	fx.Supply(config.ServiceConfig{Name: "test"}),
	fx.Supply(&config.DBConfig{
		DBHost:     "localhost",
		DBUser:     "gotemplate",
		DBName:     "gotemplate",
		DBPassword: "gotemplate",
	}),
	fx.Provide(
		common.NewLogger,
	),
	domain.Module,
	sql.Module,
	factories.Module,
)
