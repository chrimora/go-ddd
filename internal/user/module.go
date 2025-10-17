package user

import (
	"goddd/internal/common/application"
	"goddd/internal/common/infrastructure/sql"
	"goddd/internal/common/interfaces/rest"
	"goddd/internal/user/application/commands"
	"goddd/internal/user/application/eventhandlers"
	"goddd/internal/user/application/queries"
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
		queries.NewGetUserQuery,
		commands.NewCreateUserCommand,
		commands.NewUpdateUserCommand,
	),
)

var ServerModule = fx.Module(
	"user_server",
	CoreModule,
	fx.Provide(
		commonrest.AsRouteCollection(rest.NewUserRoutes),
	),
)

var ConsumerModule = fx.Module(
	"user_consumer",
	CoreModule,
	commonapplication.AsEventHandler(eventhandlers.NewUserCreatedHandler),
	commonapplication.AsEventHandler(eventhandlers.NewUserCreatedHandler2),
)

func NewUserSql(db commonsql.DBTX) *sql.Queries {
	return sql.New(db)
}
