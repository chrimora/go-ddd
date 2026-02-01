package commondomain

import (
	"context"
	"goddd/internal/common/infrastructure/sql"

	"github.com/jackc/pgx/v5"
)

type txFunc = func(tx pgx.Tx) error

type TxManager interface {
	WithTx(context.Context, txFunc) error
}

type transactionManager struct {
	db commonsql.DBTX
}

func NewTransactionManager(db commonsql.DBTX) *transactionManager {
	return &transactionManager{db: db}
}
func (m *transactionManager) WithTx(ctx context.Context, fn txFunc) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

type mockTxManager struct{}

func NewMockTxManager() *mockTxManager {
	return &mockTxManager{}
}
func (m *mockTxManager) WithTx(ctx context.Context, fn txFunc) error {
	return fn(nil)
}
