package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	outboxdb "goddd/internal/outbox/infrastructure/sql/codegen"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type OutboxRepositoryI interface {
	CreateMany(context.Context, ...domain.DomainEventI) error
	GetNextEvent(context.Context) (*outboxdb.EventOutbox, error)
	RequeueEvent(context.Context, uuid.UUID) error
	CompleteEvent(context.Context, uuid.UUID) error
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
	traceContext := infrastructure.NewTraceCtx(ctx)
	err := traceContext.IsComplete()
	if err != nil {
		return err
	}
	contextPayload, err := json.Marshal(traceContext)
	if err != nil {
		return err
	}

	t := time.Now().UTC()
	outboxSql := domain.WithTxFromCtx(e.outboxSql, ctx)
	for _, event := range events {
		eventPayload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = outboxSql.CreateEvent(
			ctx,
			outboxdb.CreateEventParams{
				ID:            uuid.New(),
				AggregateID:   event.GetAggregateId(),
				AggregateType: event.GetAggregateType(),
				EventContext:  contextPayload,
				EventType:     string(event.GetEventType()),
				Payload:       eventPayload,
				CreatedAt:     t,
				UpdatedAt:     t,
				Retries:       0,
				Status:        outboxdb.EventStatusPending,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *OutboxRepository) GetNextEvent(ctx context.Context) (*outboxdb.EventOutbox, error) {
	event, err := e.outboxSql.ClaimNextEvent(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrNotFound
		default:
			return nil, err
		}
	}
	return &event, nil
}

func (e *OutboxRepository) RequeueEvent(ctx context.Context, id uuid.UUID) error {
	return domain.WithTxFromCtx(e.outboxSql, ctx).RequeueEvent(ctx, id)
}

func (e *OutboxRepository) CompleteEvent(ctx context.Context, id uuid.UUID) error {
	return domain.WithTxFromCtx(e.outboxSql, ctx).CompleteEvent(ctx, id)
}
