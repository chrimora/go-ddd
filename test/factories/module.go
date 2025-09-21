package factories

import "go.uber.org/fx"

var Module = fx.Module(
	"factories",
	fx.Provide(
		NewUserFactory,
	),
)
