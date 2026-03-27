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

type UpdatePostTitleInput struct {
	Id    uuid.UUID
	Title string
}

type UpdatePostTitleCommand commonapplication.CommandI[UpdatePostTitleInput]

func NewUpdatePostTitleCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	postRepo domain.PostRepositoryI,
) UpdatePostTitleCommand {
	return commonapplication.NewCommand(log, &updatePostTitle{
		txManager: txManager,
		postRepo:  postRepo,
	})
}

type updatePostTitle struct {
	txManager commondomain.TxManager
	postRepo  domain.PostRepositoryI
}

func (c *updatePostTitle) Handle(
	ctx context.Context, log *slog.Logger, input UpdatePostTitleInput,
) (uuid.UUID, error) {
	post, err := c.postRepo.Get(ctx, input.Id)
	if err != nil {
		return uuid.Nil, err
	}

	post.UpdateTitle(input.Title)
	log.InfoContext(ctx, "Updating", "post", post)

	err = c.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return c.postRepo.Update(ctx, tx, post)
	})
	return post.ID(), err
}
