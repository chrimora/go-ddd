package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload" // Loads .env
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
)
var CoreModule = fx.Module(
	"config_core",
	fx.Provide(
		NewConfig[ServiceConfig],
		NewConfig[DBConfig],
	),
)

var ServerModule = fx.Module(
	"config_server",
	CoreModule,
)

var ConsumerModule = fx.Module(
	"config_consumer",
	CoreModule,
	fx.Provide(
		NewConfig[ForwarderConfig],
		NewConfig[RouterConfig],
	),
)

func NewConfig[T any]() *T {
	cfg := env.Must(env.ParseAsWithOptions[T](
		env.Options{RequiredIfNoDef: true}),
	)
	return &cfg
}
