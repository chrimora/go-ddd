package user

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		fx.Annotate(NewUserService, fx.As(new(UserServiceI))),
		fx.Annotate(NewUserRepository, fx.As(new(UserRepositoryI))),
	),
)
