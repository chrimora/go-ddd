package post

import (
	"goddd/internal/common/application"
	"goddd/internal/common/interfaces/rest"
	"goddd/internal/post/application/commands"
	"goddd/internal/post/application/eventhandlers"
	"goddd/internal/post/application/queries"
	"goddd/internal/post/domain"
	"goddd/internal/post/infrastructure/sql"
	"goddd/internal/post/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"post_core",
	fx.Provide(
		sql.NewWritePostSql,
		sql.NewReadPostSql,
		fx.Annotate(domain.NewPostRepository, fx.As(new(domain.PostRepositoryI))),
		queries.NewGetPostQuery,
		queries.NewGetPostsQuery,
		commands.NewCreatePostCommand,
		commands.NewUpdatePostTitleCommand,
	),
)

var APIModule = fx.Module(
	"post_api",
	CoreModule,
	fx.Provide(
		commonrest.AsRouteCollection(rest.NewPostRoutes),
	),
)

var ConsumerModule = fx.Module(
	"post_consumer",
	CoreModule,
	commonapplication.AsEventHandler(eventhandlers.NewPostCreatedHandler),
)
