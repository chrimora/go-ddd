package common

import "go.uber.org/fx"

var Module = fx.Module(
	"common",
	fx.Provide(
		NewLogger,
		NewContext,
		NewDBPool,
	),
)
