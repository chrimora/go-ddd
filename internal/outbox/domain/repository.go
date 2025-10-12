package domain

import (
	"context"
	"encoding/json"
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	outboxsql "goddd/internal/outbox/infrastructure/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type OutboxRepositoryI interface {
	CreateMany(context.Context, ...commondomain.DomainEventI) error
	GetNextEventBatch(context.Context, int, int) ([]*OutboxEvent, error)
	RequeueStaleEvents(context.Context, time.Time, int) (int, error)
	CompleteEvent(context.Context, uuid.UUID) error
}

type OutboxRepository OutboxRepositoryI
type outboxRepository struct {
	log       *slog.Logger
	outboxSql *outboxsql.Queries
}

func NewOutboxRepository(log *slog.Logger, outboxSql *outboxsql.Queries) OutboxRepository {
	return &outboxRepository{
		log:       log,
		outboxSql: outboxSql,
	}
}

func (e *outboxRepository) CreateMany(
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

func (e *outboxRepository) GetNextEventBatch(
	ctx context.Context, batchSize, retries int,
) ([]*OutboxEvent, error) {
	events, err := e.outboxSql.ClaimNextEventBatch(
		ctx,
		outboxsql.ClaimNextEventBatchParams{Limit: int32(batchSize), Retries: int32(retries)},
	)
	if err != nil {
		return nil, err
	}

	ret := make([]*OutboxEvent, len(events))
	for i, e := range events {
		ret[i] = &OutboxEvent{
			ID:            e.ID,
			AggregateID:   e.AggregateID,
			AggregateType: e.AggregateType,
			EventContext:  e.EventContext,
			EventType:     e.EventType,
			Payload:       e.Payload,
			CreatedAt:     e.CreatedAt,
			UpdatedAt:     e.UpdatedAt,
		}
	}
	return ret, nil
}

func (e *outboxRepository) RequeueStaleEvents(
	ctx context.Context, before time.Time, retries int,
) (int, error) {
	ids, err := commondomain.WithTxFromCtx(e.outboxSql, ctx).RequeueStaleEvents(
		ctx,
		outboxsql.RequeueStaleEventsParams{UpdatedAt: before, Retries: int32(retries)},
	)
	return len(ids), err
}

func (e *outboxRepository) CompleteEvent(ctx context.Context, id uuid.UUID) error {
	return commondomain.WithTxFromCtx(e.outboxSql, ctx).CompleteEvent(ctx, id)
}
