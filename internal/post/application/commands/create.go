package commands

import (
	"context"
	commonapplication "goddd/internal/common/application"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/post/domain"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreatePostInput struct {
	Title  string
	Author string
}

type CreatePostCommand commonapplication.CommandI[CreatePostInput]

func NewCreatePostCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	postRepo domain.PostRepositoryI,
) CreatePostCommand {
	return commonapplication.NewCommand(log, &createPost{
		txManager: txManager,
		postRepo:  postRepo,
	})
}

type createPost struct {
	txManager commondomain.TxManager
	postRepo  domain.PostRepositoryI
}

func (c *createPost) Handle(
	ctx context.Context, log *slog.Logger, input CreatePostInput,
) (uuid.UUID, error) {
	post := domain.NewPost(input.Title, input.Author)
	log.InfoContext(ctx, "Creating", "post", post)

	err := c.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return c.postRepo.Create(ctx, tx, post)
	})
	return post.ID(), err
}
