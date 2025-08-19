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
		NewUserSql,
	),
)

func NewUserSql(db common.DBTX) UserSqlI {
	return userdb.New(db)
}
