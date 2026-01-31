package test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/common/test"
	"goddd/internal/user/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserFactory commontest.Mock[domain.User]
type userFactory struct {
	repo      domain.UserRepositoryI
	txManager *commondomain.TxManager
}

func NewUserFactory(repo domain.UserRepositoryI, txManager *commondomain.TxManager) UserFactory {
	return &userFactory{repo: repo, txManager: txManager}
}

func (f *userFactory) Mock(t *testing.T, overrides ...map[string]any) *domain.User {
	ctx := context.Background()
	fields := &struct {
		ID        uuid.UUID
		Version   int
		CreatedAt time.Time
		UpdatedAt time.Time
		Name      string
	}{
		ID:   uuid.New(),
		Name: "Christopher",
	}
	commontest.Merge(fields, overrides)

	user := domain.RehydrateUser(fields.ID, fields.Version, fields.CreatedAt, fields.UpdatedAt, fields.Name)
	f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return f.repo.Create(ctx, tx, user)
	})

	t.Cleanup(func() {
		f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
			return f.repo.Remove(ctx, tx, user)
		})
	})
	return user
}
