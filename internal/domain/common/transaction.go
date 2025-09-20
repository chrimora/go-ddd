package common

import (
	"context"
	"gotemplate/internal/common"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Tx struct {
	log              *slog.Logger
	outboxRepository OutboxRepositoryI

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
	var events []DomainEventI
	for _, agg := range tx.aggregates {
		events = append(events, agg.PullEvents()...)
	}

	for _, event := range events {
		serviceContext, err := common.NewServiceCtx(tx.ctx)
		if err != nil {
			return err
		}
		err = tx.outboxRepository.WithTx(tx.Tx).Create(tx.ctx, serviceContext, event)
		if err != nil {
			return err
		}

	}
	return tx.Tx.Commit(tx.ctx)
}

type TxFactory func(context.Context) (*Tx, error)

func NewTxFactory(db common.DBTX, log *slog.Logger, outboxRepository OutboxRepositoryI) TxFactory {
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
