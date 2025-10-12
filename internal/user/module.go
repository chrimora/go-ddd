package user

import (
	"goddd/internal/common/application"
	"goddd/internal/common/infrastructure/sql"
	"goddd/internal/common/interfaces/rest"
	"goddd/internal/user/application"
	"goddd/internal/user/application/eventhandlers"
	"goddd/internal/user/domain"
	"goddd/internal/user/infrastructure/sql"
	"goddd/internal/user/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"user_core",
	fx.Provide(
		NewUserSql,
		fx.Annotate(domain.NewUserRepository, fx.As(new(domain.UserRepositoryI))),
		fx.Annotate(application.NewUserService, fx.As(new(application.UserServiceI))),
	),
)

var ServerModule = fx.Module(
	"user_server",
	CoreModule,
	fx.Provide(
		commonrest.AsRouteCollection(rest.NewUserRoutes),
	),
)

var WorkerModule = fx.Module(
	"user_worker",
	CoreModule,
	commonapplication.AsEventHandler(eventhandlers.NewUserCreatedHandler),
	commonapplication.AsEventHandler(eventhandlers.NewUserCreatedHandler2),
)

func NewUserSql(db commonsql.DBTX) *sql.Queries {
	return sql.New(db)
}
