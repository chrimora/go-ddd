package common

import (
	"context"
	"goddd/internal/common"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Tx struct {
	log              *slog.Logger
	outboxRepository OutboxRepositoryI

	Tx     pgx.Tx
	ctx    context.Context
	events []DomainEventI
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback(tx.ctx)
}

func (tx *Tx) AddEvents(events ...DomainEventI) {
	tx.events = append(tx.events, events...)
}

func (tx *Tx) Commit() error {
	for _, event := range tx.events {
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
