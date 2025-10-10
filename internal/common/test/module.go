package test

import (
	"goddd/internal/common"
	"goddd/internal/config"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"test_core",
	fx.Supply(config.ServiceConfig{Name: "test"}),
	fx.Supply(&config.DBConfig{
		DBHost:     "localhost",
		DBUser:     "goddd",
		DBName:     "goddd",
		DBPassword: "goddd",
	}),
	common.Module,
)
