package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/post/domain"
	"log/slog"

	"github.com/google/uuid"
)

type GetPostInput struct {
	Id uuid.UUID
}

type GetPostQuery commonapplication.QueryI[GetPostInput, *domain.Post]

func NewGetPostQuery(
	log *slog.Logger,
	postRepo domain.PostRepositoryI,
) GetPostQuery {
	return commonapplication.NewQuery(log, &getPost{
		postRepo: postRepo,
	})
}

type getPost struct {
	postRepo domain.PostRepositoryI
}

func (q *getPost) Handle(
	ctx context.Context, log *slog.Logger, input GetPostInput,
) (*domain.Post, error) {
	post, err := q.postRepo.Get(ctx, input.Id)
	log.InfoContext(ctx, "Got", "post", post)
	return post, err
}
