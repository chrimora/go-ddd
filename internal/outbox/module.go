package outbox

import (
	"goddd/internal/common/infrastructure/sql"
	"goddd/internal/outbox/application"
	"goddd/internal/outbox/domain"
	"goddd/internal/outbox/infrastructure/pubsub"
	"goddd/internal/outbox/infrastructure/sql"

	"github.com/ThreeDotsLabs/watermill/message"
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

var ConsumerModule = fx.Module(
	"outbox_consumer",
	CoreModule,
	fx.Provide(
		fx.Annotate(pubsub.NewGoChannel, fx.As(new(message.Publisher)), fx.As(new(message.Subscriber))),
		fx.Annotate(application.NewEventPublisher, fx.As(new(application.EventPublisherI))),
		pubsub.NewRouter,
		application.NewForwarder,
		fx.Annotate(application.NewConsumer, fx.ParamTags(`group:"eventHandlers"`)),
	),
)

func NewOutboxSql(db commonsql.DBTX) *sql.Queries {
	return sql.New(db)
}
