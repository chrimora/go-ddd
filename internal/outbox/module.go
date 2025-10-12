package outbox

import (
	"goddd/internal/common/infrastructure/sql"
	"goddd/internal/outbox/application"
	"goddd/internal/outbox/domain"
	"goddd/internal/outbox/infrastructure/sql"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"outbox_core",
	fx.Provide(
		NewOutboxSql,
		fx.Annotate(domain.NewOutboxRepository, fx.As(new(OutboxRepositoryI))),
	),
)

var ServerModule = fx.Module(
	"outbox_server",
	CoreModule,
)

var WorkerModule = fx.Module(
	"outbox_worker",
	CoreModule,
	fx.Provide(
		fx.Annotate(application.NewDispatcher, fx.ParamTags(`group:"eventHandlers"`)),
		application.NewWorker,
	),
)

func NewOutboxSql(db commonsql.DBTX) *sql.Queries {
	return sql.New(db)
}
