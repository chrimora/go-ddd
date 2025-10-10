package outbox

import (
	"goddd/internal/common/infrastructure/sql"
	"goddd/internal/outbox/application"
	"goddd/internal/outbox/domain"
	outboxdb "goddd/internal/outbox/infrastructure/sql/codegen"

	"go.uber.org/fx"
)

var CoreModule = fx.Module(
	"outbox_core",
	fx.Provide(
		NewOutboxSql,
		fx.Annotate(domain.NewOutboxRepository, fx.As(new(domain.OutboxRepositoryI))),
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

func NewOutboxSql(db sql.DBTX) *outboxdb.Queries {
	return outboxdb.New(db)
}
