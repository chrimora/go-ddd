package domain

import (
	"context"
	"goddd/internal/common/infrastructure/sql"

	"github.com/jackc/pgx/v5"
)

var txKey = "tx"

type withTxFunc[T any] interface {
	WithTx(tx pgx.Tx) T
}

// Return sql with a transaction from the context (if it exists)
func WithTxFromCtx[T withTxFunc[T]](sql T, ctx context.Context) T {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok && tx != nil {
		return sql.WithTx(tx)
	}
	return sql
}

type TxFactory func(context.Context) (pgx.Tx, error)

func NewTxFactory(db sql.DBTX) TxFactory {
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

type TxFunc func(ctx context.Context) error

// Provide a context with a new transaction
func (m *TxManager) WithTxCtx(ctx context.Context, fn TxFunc) error {
	tx, err := m.txFactory(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
