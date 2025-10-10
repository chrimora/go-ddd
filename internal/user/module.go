package user

import (
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure/sql"
	rest "goddd/internal/common/interfaces/rest"
	user_application "goddd/internal/user/application"
	eventhandlers "goddd/internal/user/application/event_handlers"
	user_domain "goddd/internal/user/domain"
	userdb "goddd/internal/user/infrastructure/sql/codegen"
	user_rest "goddd/internal/user/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"user_core",
	fx.Provide(
		NewUserSql,
		fx.Annotate(user_domain.NewUserRepository, fx.As(new(user_application.UserRepositoryI))),
		fx.Annotate(user_application.NewUserService, fx.As(new(user_application.UserServiceI))),
	),
)

var ServerModule = fx.Module(
	"user_server",
	CoreModule,
	fx.Provide(
		rest.AsRouteCollection(user_rest.NewUserRoutes),
	),
)

var WorkerModule = fx.Module(
	"user_worker",
	CoreModule,
	domain.AsEventHandler(eventhandlers.NewUserCreatedHandler, &eventhandlers.UserCreatedHandler{}),
	domain.AsEventHandler(eventhandlers.NewUserCreatedHandler2, &eventhandlers.UserCreatedHandler2{}),
)

func NewUserSql(db sql.DBTX) *userdb.Queries {
	return userdb.New(db)
}
