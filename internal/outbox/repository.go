package outbox

import (
	"context"
	"encoding/json"
	"gotemplate/internal/common"
	domain "gotemplate/internal/domain/common"
	outboxdb "gotemplate/internal/outbox/db"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type OutboxRepository struct {
	log       *slog.Logger
	outboxSql *outboxdb.Queries
}

func NewOutboxRepository(log *slog.Logger, outboxSql *outboxdb.Queries) domain.OutboxRepositoryI {
	return &OutboxRepository{
		log:       log,
		outboxSql: outboxSql,
	}
}

func (e *OutboxRepository) WithTx(tx pgx.Tx) domain.OutboxRepositoryI {
	return NewOutboxRepository(e.log, e.outboxSql.WithTx(tx))
}

func (e *OutboxRepository) Create(
	ctx context.Context,
	eventContext common.ServiceContext,
	event domain.DomainEventI,
) error {
	_eventContext, err := json.Marshal(eventContext)
	if err != nil {
		return err
	}
	_event, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return e.outboxSql.CreateEvent(
		ctx,
		outboxdb.CreateEventParams{
			AggregateID:  event.GetAggregateId(),
			EventContext: _eventContext,
			EventType:    string(event.GetType()),
			Payload:      _event,
		},
	)
}
