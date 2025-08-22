package user

import (
	"gotemplate/internal/common"
	userdb "gotemplate/internal/domain/user/db"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		common.AsRouteCollection(NewUserRoutes),
		NewUserService,
		NewUserRepository,
		NewSql,
	),
)

func NewSql(db common.DBTX) *userdb.Queries {
	return userdb.New(db)
}
