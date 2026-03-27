package test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	commontest "goddd/internal/common/test"
	"goddd/internal/post/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostFactory commontest.Mock[domain.Post]
type postFactory struct {
	repo      domain.PostRepositoryI
	txManager commondomain.TxManager
}

func NewPostFactory(repo domain.PostRepositoryI, txManager commondomain.TxManager) PostFactory {
	return &postFactory{repo: repo, txManager: txManager}
}

func (f *postFactory) Mock(t *testing.T, ctx context.Context, overrides ...map[string]any) *domain.Post {
	fields := &struct {
		ID          uuid.UUID
		Version     int
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Title       string
		PublishDate time.Time
		Author      string
	}{
		ID:          commondomain.NewUUID(),
		Title:       "A Post Title",
		PublishDate: time.Now().UTC(),
		Author:      "Anonymous",
	}
	commontest.Merge(fields, overrides)

	post := domain.RehydratePost(
		fields.ID, fields.Version, fields.CreatedAt, fields.UpdatedAt,
		fields.Title, fields.PublishDate, fields.Author,
	)
	err := f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return f.repo.Create(ctx, tx, post)
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		err := f.txManager.WithTx(ctx, func(tx pgx.Tx) error {
			return f.repo.Remove(ctx, tx, post)
		})
		if err != nil {
			panic(err)
		}
	})
	return post
}
