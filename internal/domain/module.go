package domain

import (
	"gotemplate/internal/domain/common"
	"gotemplate/internal/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"domain",
	fx.Provide(
		common.NewTxFactory,
		common.NewOutboxRepository,
	),
	user.Module,
)

type (
	User         = user.User
	UserServiceI = user.UserServiceI
)
