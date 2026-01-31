package commondomain

import (
	"context"
	"goddd/internal/common/infrastructure/sql"

	"github.com/jackc/pgx/v5"
)

type TxManager struct {
	db commonsql.DBTX
}

func NewTxManager(db commonsql.DBTX) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
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
