package user

import (
	"gotemplate/internal/common"
	userdb "gotemplate/internal/infrastructure/sql/codegen/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewUserService,
		NewUserRepository,
		NewSql,
	),
)

func NewSql(db common.DBTX) *userdb.Queries {
	return userdb.New(db)
}
