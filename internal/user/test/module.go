package test

import (
	commoninfrastructure "goddd/internal/common/infrastructure"
	"goddd/internal/common/test"
	"goddd/internal/user"
	"goddd/internal/user/application/commands"
	"goddd/internal/user/domain"

	"go.uber.org/fx"
)

var UnitTestModule = fx.Module(
	"user_unit_test",
	commontest.UnitTestModule,
	fx.Provide(
		fx.Annotate(
			commoninfrastructure.NewInMemoryRepository[*domain.User],
			fx.As(new(domain.UserRepositoryI)),
		),
		NewUserFactory,
		commands.NewCreateUserCommand,
		commands.NewUserChangeNameCommand,
	),
)

var IntegrationTestModule = fx.Module(
	"user_integration_test",
	commontest.IntegrationTestModule,
	user.CoreModule,
	fx.Provide(
		NewUserFactory,
	),
)
