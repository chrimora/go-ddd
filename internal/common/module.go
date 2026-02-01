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
		commoninfrastructure.NewLogger,
		commonsql.NewContext,
		commonsql.NewDBPool,
		fx.Annotate(commondomain.NewTransactionManager, fx.As(new(commondomain.TxManager))),
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
