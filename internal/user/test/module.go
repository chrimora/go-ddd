package test

import (
	"goddd/internal/common/test"
	"goddd/internal/user"

	"go.uber.org/fx"
)

var IntegrationTestModule = fx.Module(
	"user_test",
	test.CoreModule,
	user.CoreModule,
	fx.Provide(
		NewUserFactory,
	),
)
