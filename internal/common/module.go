package common

import (
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	"goddd/internal/common/infrastructure/sql"
	commonrest "goddd/internal/common/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"common_core",
	fx.Provide(
		infrastructure.NewLogger,
		commondomain.NewTxFactory,
		commondomain.NewTxManager,
		commonsql.NewContext,
		commonsql.NewDBPool,
	),
)

var ServerModule = fx.Module(
	"common_server",
	CoreModule,
	fx.Provide(
		commonrest.NewHTTPServer,
		commonrest.NewServeMux,
		fx.Annotate(
			commonrest.NewApi,
			fx.ParamTags(`group:"routeCollection"`),
		),
	),
)

var ConsumerModule = fx.Module(
	"common_consumer",
	CoreModule,
)
