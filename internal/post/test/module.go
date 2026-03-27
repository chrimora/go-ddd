package test

import (
	commoninfrastructure "goddd/internal/common/infrastructure"
	commontest "goddd/internal/common/test"
	"goddd/internal/post"
	"goddd/internal/post/application/commands"
	"goddd/internal/post/domain"

	"go.uber.org/fx"
)

var UnitTestModule = fx.Module(
	"post_unit_test",
	commontest.UnitTestModule,
	fx.Provide(
		fx.Annotate(
			commoninfrastructure.NewInMemoryRepository[*domain.Post],
			fx.As(new(domain.PostRepositoryI)),
		),
		NewPostFactory,
		commands.NewCreatePostCommand,
		commands.NewUpdatePostTitleCommand,
	),
)

var IntegrationTestModule = fx.Module(
	"post_integration_test",
	commontest.IntegrationTestModule,
	post.CoreModule,
	fx.Provide(
		NewPostFactory,
	),
)
