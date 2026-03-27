package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/post/domain"
	postsql "goddd/internal/post/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
)

type GetPostsInput struct {
	Limit int
	After *uuid.UUID
	Title *string
}

type GetPostsOutput struct {
	Posts []domain.Post
	Next  *uuid.UUID
}

type GetPostsQuery commonapplication.QueryI[GetPostsInput, GetPostsOutput]

func NewGetPostsQuery(
	log *slog.Logger,
	postSql postsql.ReadPostSql,
) GetPostsQuery {
	return commonapplication.NewQuery(log, &getPosts{
		postSql: postSql,
	})
}

type getPosts struct {
	postSql postsql.ReadPostSql
}

func (q *getPosts) Handle(
	ctx context.Context, log *slog.Logger, input GetPostsInput,
) (GetPostsOutput, error) {
	rows, err := q.postSql.ListPosts(
		ctx,
		postsql.ListPostsParams{
			LimitPlusOne: int32(input.Limit + 1),
			After:        input.After,
			Title:        input.Title,
		},
	)
	if err != nil {
		return GetPostsOutput{}, err
	}

	var next *uuid.UUID
	if len(rows) > input.Limit {
		rows = rows[:input.Limit]
		next = &rows[input.Limit-1].ID
	}

	posts := make([]domain.Post, len(rows))
	for i, row := range rows {
		posts[i] = *domain.RehydratePost(
			row.ID, int(row.Version), row.CreatedAt, row.UpdatedAt,
			row.Title, row.PublishDate, row.Author,
		)
	}
	return GetPostsOutput{Posts: posts, Next: next}, nil
}
