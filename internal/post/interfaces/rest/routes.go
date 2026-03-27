package rest

import (
	"context"
	"errors"
	commondomain "goddd/internal/common/domain"
	commonrest "goddd/internal/common/interfaces/rest"
	"goddd/internal/post/application/commands"
	"goddd/internal/post/application/queries"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type (
	PostRoutes commonrest.RouteCollection
	postRoutes struct {
		log             *slog.Logger
		getPost         queries.GetPostQuery
		getPosts        queries.GetPostsQuery
		createPost      commands.CreatePostCommand
		updatePostTitle commands.UpdatePostTitleCommand
	}
)

func NewPostRoutes(
	log *slog.Logger,
	getPost queries.GetPostQuery,
	getPosts queries.GetPostsQuery,
	createPost commands.CreatePostCommand,
	updatePostTitle commands.UpdatePostTitleCommand,
) PostRoutes {
	return &postRoutes{
		log:             log,
		getPost:         getPost,
		getPosts:        getPosts,
		createPost:      createPost,
		updatePostTitle: updatePostTitle,
	}
}

func (r *postRoutes) Register(api huma.API) {
	huma.Get(api, "/posts/{id}", r.get)
	huma.Get(api, "/posts", r.list)
	huma.Post(api, "/posts", r.create)
	huma.Put(api, "/posts/{id}/update-title", r.updateTitle)
}

type PostProfile struct {
	commonrest.IdPayload
	Title       string    `json:"title"`
	PublishDate time.Time `json:"publishDate"`
	Author      string    `json:"author"`
}

func (r *postRoutes) get(
	ctx context.Context, req *commonrest.IdParam,
) (*commonrest.Response[PostProfile], error) {
	post, err := r.getPost.Handle(ctx, queries.GetPostInput{Id: req.ID})
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(r.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(r.log, ctx, err)
		}
	}
	res := PostProfile{}
	res.ID = post.ID()
	res.Title = post.Title()
	res.PublishDate = post.PublishDate()
	res.Author = post.Author()
	return commonrest.BuildResponse(res), nil
}

type PostsQuery struct {
	commonrest.PaginationQuery
	Title string `query:"title"`
}

func (r *postRoutes) list(
	ctx context.Context, req *PostsQuery,
) (*commonrest.Response[commonrest.Page[PostProfile]], error) {
	var after *uuid.UUID
	if req.AfterCursor != uuid.Nil {
		after = &req.AfterCursor
	}
	var title *string
	if req.Title != "" {
		title = &req.Title
	}

	out, err := r.getPosts.Handle(ctx, queries.GetPostsInput{
		Limit: req.Limit,
		After: after,
		Title: title,
	})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(r.log, ctx, err)
	}

	items := make([]PostProfile, len(out.Posts))
	for i, post := range out.Posts {
		p := PostProfile{}
		p.ID = post.ID()
		p.Title = post.Title()
		p.PublishDate = post.PublishDate()
		p.Author = post.Author()
		items[i] = p
	}
	res := commonrest.Page[PostProfile]{Items: items, NextCursor: out.Next}
	return commonrest.BuildResponse(res), nil
}

type CreatePostPayload struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func (r *postRoutes) create(
	ctx context.Context, req *commonrest.CreateRequest[CreatePostPayload],
) (*commonrest.Response[commonrest.IdPayload], error) {
	id, err := r.createPost.Handle(ctx, commands.CreatePostInput{
		Title:  req.Body.Title,
		Author: req.Body.Author,
	})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(r.log, ctx, err)
	}
	return commonrest.BuildResponse(commonrest.IdPayload{ID: id}), nil
}

type UpdateTitlePayload struct {
	Title string `json:"title"`
}

func (r *postRoutes) updateTitle(
	ctx context.Context, req *commonrest.UpdateRequest[UpdateTitlePayload],
) (*commonrest.EmptyResponse, error) {
	_, err := r.updatePostTitle.Handle(ctx, commands.UpdatePostTitleInput{
		Id:    req.ID,
		Title: req.Body.Title,
	})
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(r.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(r.log, ctx, err)
		}
	}
	return &commonrest.EmptyResponse{}, nil
}
