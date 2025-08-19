package user

import (
	"gotemplate/internal/common"
	userdb "gotemplate/internal/user/db"

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
