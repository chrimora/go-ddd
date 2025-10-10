package common

import (
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	"goddd/internal/common/infrastructure/sql"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"common",
	fx.Provide(
		infrastructure.NewLogger,
		domain.NewTxFactory,
		domain.NewTxManager,
		sql.NewContext,
		sql.NewDBPool,
	),
)
