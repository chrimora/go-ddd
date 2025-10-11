package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	outboxsql "goddd/internal/outbox/infrastructure/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type OutboxRepositoryI interface {
	CreateMany(context.Context, ...commondomain.DomainEventI) error
	GetNextEvent(context.Context) (*outboxsql.EventOutbox, error)
	RequeueEvent(context.Context, uuid.UUID) error
	CompleteEvent(context.Context, uuid.UUID) error
}

type OutboxRepository struct {
	log       *slog.Logger
	outboxSql *outboxsql.Queries
}

func NewOutboxRepository(log *slog.Logger, outboxSql *outboxsql.Queries) OutboxRepositoryI {
	return &OutboxRepository{
		log:       log,
		outboxSql: outboxSql,
	}
}

func (e *OutboxRepository) CreateMany(
	ctx context.Context,
	events ...commondomain.DomainEventI,
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
	outboxQ := commondomain.WithTxFromCtx(e.outboxSql, ctx)
	for _, event := range events {
		eventPayload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = outboxQ.CreateEvent(
			ctx,
			outboxsql.CreateEventParams{
				ID:            uuid.New(),
				AggregateID:   event.GetAggregateId(),
				AggregateType: event.GetAggregateType(),
				EventContext:  contextPayload,
				EventType:     string(event.GetEventType()),
				Payload:       eventPayload,
				CreatedAt:     t,
				UpdatedAt:     t,
				Retries:       0,
				Status:        outboxsql.EventStatusPending,
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *OutboxRepository) GetNextEvent(ctx context.Context) (*outboxsql.EventOutbox, error) {
	event, err := e.outboxSql.ClaimNextEvent(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, commondomain.ErrNotFound
		default:
			return nil, err
		}
	}
	return &event, nil
}

func (e *OutboxRepository) RequeueEvent(ctx context.Context, id uuid.UUID) error {
	return commondomain.WithTxFromCtx(e.outboxSql, ctx).RequeueEvent(ctx, id)
}

func (e *OutboxRepository) CompleteEvent(ctx context.Context, id uuid.UUID) error {
	return commondomain.WithTxFromCtx(e.outboxSql, ctx).CompleteEvent(ctx, id)
}
