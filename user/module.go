package user

import (
	"gotemplate/common"
	userdb "gotemplate/user/db"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		common.AsRouteCollection(NewUserRoutes),
		NewUserService,
		NewUserRepository,
	),
)

func NewUserRepository(db common.DBTX) UserRepositoryI {
	return userdb.New(db)
}
