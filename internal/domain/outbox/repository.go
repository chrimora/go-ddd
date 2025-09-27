package outbox

import (
	"context"
	"encoding/json"
	"goddd/internal/common"
	domain "goddd/internal/domain/common"
	outboxdb "goddd/internal/infrastructure/sql/codegen/outbox"
	"log/slog"
)

type OutboxRepositoryI interface {
	CreateMany(context.Context, ...domain.DomainEventI) error
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

func (e *OutboxRepository) CreateMany(
	ctx context.Context,
	events ...domain.DomainEventI,
) error {
	traceContext, err := common.NewTraceCtx(ctx)
	if err != nil {
		return err
	}
	contextPayload, err := json.Marshal(traceContext)
	if err != nil {
		return err
	}

	for _, event := range events {
		eventPayload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = domain.WithTxFromCtx(e.outboxSql, ctx).CreateEvent(
			ctx,
			outboxdb.CreateEventParams{
				AggregateID:  event.GetAggregateId(),
				EventContext: contextPayload,
				EventType:    string(event.GetType()),
				Payload:      eventPayload,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
