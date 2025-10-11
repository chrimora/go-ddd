package user

import (
	"goddd/internal/common/domain"
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
		fx.Annotate(domain.NewUserRepository, fx.As(new(application.UserRepositoryI))),
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
	commondomain.AsEventHandler(eventhandlers.NewUserCreatedHandler, &eventhandlers.UserCreatedHandler{}),
	commondomain.AsEventHandler(eventhandlers.NewUserCreatedHandler2, &eventhandlers.UserCreatedHandler2{}),
)

func NewUserSql(db commonsql.DBTX) *sql.Queries {
	return sql.New(db)
}
