package commondomain

import (
	"context"
	"goddd/internal/common/infrastructure/sql"

	"github.com/jackc/pgx/v5"
)

type TxFactory func(context.Context) (pgx.Tx, error)

func NewTxFactory(db commonsql.DBTX) TxFactory {
	return func(ctx context.Context) (pgx.Tx, error) {
		return db.Begin(ctx)
	}
}

type TxManager struct {
	txFactory TxFactory
}

func NewTxManager(txFactory TxFactory) *TxManager {
	return &TxManager{txFactory: txFactory}
}

type TxFunc func(tx pgx.Tx) error

func (m *TxManager) WithTx(ctx context.Context, fn TxFunc) error {
	tx, err := m.txFactory(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
