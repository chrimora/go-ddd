package common

import (
	"context"
	"gotemplate/internal/common"
	"gotemplate/internal/outbox"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type AggregateRootI interface {
	PullEvents() []outbox.DomainEventI
}

type Tx struct {
	log              *slog.Logger
	outboxRepository outbox.OutboxRepositoryI

	Tx         pgx.Tx
	ctx        context.Context
	aggregates []AggregateRootI
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback(tx.ctx)
}

func (tx *Tx) TrackEvents(aggregate AggregateRootI) {
	tx.aggregates = append(tx.aggregates, aggregate)
}

func (tx *Tx) Commit() error {
	var events []outbox.DomainEventI
	for _, agg := range tx.aggregates {
		events = append(events, agg.PullEvents()...)
	}

	for _, event := range events {
		payload, err := outbox.CreateEventPayload(tx.ctx, event)
		if err != nil {
			return err
		}
		err = tx.outboxRepository.WithTx(tx.Tx).Create(tx.ctx, event.Type(), payload)
		if err != nil {
			return err
		}

	}
	return tx.Tx.Commit(tx.ctx)
}

type TxFactory func(context.Context) (*Tx, error)

func NewTxFactory(db common.DBTX, log *slog.Logger, outboxRepository outbox.OutboxRepositoryI) TxFactory {
	return func(ctx context.Context) (*Tx, error) {
		tx, err := db.Begin(ctx)
		return &Tx{
			log:              log,
			outboxRepository: outboxRepository,
			Tx:               tx,
			ctx:              ctx,
		}, err
	}
}
