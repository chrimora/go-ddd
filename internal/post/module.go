package post

import (
	commonapplication "goddd/internal/common/application"
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/post/application/commands"
	"goddd/internal/post/application/eventhandlers"
	"goddd/internal/post/application/queries"
	"goddd/internal/post/domain"
	postsql "goddd/internal/post/infrastructure/sql"
	postrest "goddd/internal/post/interfaces/rest"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"post_core",
	fx.Provide(
		postsql.NewWritePostSql,
		postsql.NewReadPostSql,
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
		commonrest.AsRouteCollection(postrest.NewPostRoutes),
	),
)

var ConsumerModule = fx.Module(
	"post_consumer",
	CoreModule,
	commonapplication.AsEventHandler(eventhandlers.NewPostCreatedHandler),
)
