package outbox

import (
	"context"
	outboxdb "gotemplate/internal/outbox/db"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type OutboxRepositoryI interface {
	WithTx(pgx.Tx) OutboxRepositoryI
	Create(context.Context, EventType, []byte) error
}

type OutboxRepository struct {
	log       *slog.Logger
	outboxSql *outboxdb.Queries
}

func NewOutboxRepository(log *slog.Logger, outboxSql *outboxdb.Queries) OutboxRepositoryI {
	return &OutboxRepository{
		log:       log,
		outboxSql: outboxSql,
	}
}

func (e *OutboxRepository) WithTx(tx pgx.Tx) OutboxRepositoryI {
	return NewOutboxRepository(e.log, e.outboxSql.WithTx(tx))
}

func (e *OutboxRepository) Create(ctx context.Context, eventType EventType, payload []byte) error {
	return e.outboxSql.CreateEvent(
		ctx,
		outboxdb.CreateEventParams{EventType: string(eventType), Payload: payload},
	)
}
