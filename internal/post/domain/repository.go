package domain

import (
	"context"
	"database/sql"
	"errors"
	"goddd/internal/outbox"
	postsql "goddd/internal/post/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostRepositoryI interface {
	Get(context.Context, uuid.UUID) (*Post, error)
	Create(context.Context, pgx.Tx, *Post) error
	Update(context.Context, pgx.Tx, *Post) error
	Remove(context.Context, pgx.Tx, *Post) error
}

type PostRepository PostRepositoryI

type postRepository struct {
	log        *slog.Logger
	postSql    postsql.WritePostSql
	outboxRepo outbox.OutboxRepositoryI
}

func NewPostRepository(
	log *slog.Logger,
	postSql postsql.WritePostSql,
	outboxRepo outbox.OutboxRepositoryI,
) PostRepository {
	return &postRepository{
		log:        log,
		postSql:    postSql,
		outboxRepo: outboxRepo,
	}
}

func (r *postRepository) Get(ctx context.Context, id uuid.UUID) (*Post, error) {
	row, err := r.postSql.GetPost(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound(id)
		default:
			return nil, err
		}
	}
	return RehydratePost(
		row.ID, int(row.Version), row.CreatedAt, row.UpdatedAt,
		row.Title, row.PublishDate, row.Author,
	), nil
}

func (r *postRepository) Create(ctx context.Context, tx pgx.Tx, post *Post) error {
	_, err := r.postSql.WithTx(tx).CreatePost(
		ctx,
		postsql.CreatePostParams{
			ID:          post.ID(),
			Version:     int32(post.Version()),
			CreatedAt:   post.CreatedAt(),
			UpdatedAt:   post.UpdatedAt(),
			Title:       post.Title(),
			PublishDate: post.PublishDate(),
			Author:      post.Author(),
		},
	)
	if err != nil {
		return err
	}
	return r.outboxRepo.CreateMany(ctx, tx, post.PullEvents()...)
}

func (r *postRepository) Update(ctx context.Context, tx pgx.Tx, post *Post) error {
	_, err := r.postSql.WithTx(tx).UpdatePost(
		ctx,
		postsql.UpdatePostParams{
			ID:        post.ID(),
			Version:   int32(post.Version()),
			UpdatedAt: post.UpdatedAt(),
			Title:     post.Title(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRaceCondition(post.ID())
		default:
			return err
		}
	}
	return r.outboxRepo.CreateMany(ctx, tx, post.PullEvents()...)
}

func (r *postRepository) Remove(ctx context.Context, tx pgx.Tx, post *Post) error {
	return r.postSql.WithTx(tx).RemovePost(ctx, post.ID())
}
