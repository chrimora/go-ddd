package outbox

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbox",
	fx.Provide(
		fx.Annotate(NewOutboxRepository, fx.As(new(OutboxRepositoryI))),
	),
)
