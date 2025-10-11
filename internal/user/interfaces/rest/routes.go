package rest

import (
	"context"
	"errors"
	"goddd/internal/common/domain"
	"goddd/internal/common/interfaces/rest"
	"goddd/internal/user/application"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

type UserRoutes struct {
	log         *slog.Logger
	userService application.UserServiceI
}

func NewUserRoutes(log *slog.Logger, userService application.UserServiceI) *UserRoutes {
	return &UserRoutes{
		log:         log,
		userService: userService,
	}
}

func (u *UserRoutes) Register(api huma.API) {
	huma.Get(api, "/user/{id}", u.get)
	huma.Post(api, "/user", u.create)
	huma.Put(api, "/user/{id}", u.update)
}

type UserCreatePayload struct {
	Name string `json:"name"`
}
type UserUpdatePayload struct {
	UserCreatePayload
}

type UserResponse struct {
	commonrest.IdPayload
	Name string `json:"name"`
}

func (u *UserRoutes) get(
	ctx context.Context, req *commonrest.IdParam,
) (*commonrest.Response[UserResponse], error) {
	user, err := u.userService.Get(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(u.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	return commonrest.BuildResponse(res), nil
}

func (u *UserRoutes) create(
	ctx context.Context, req *commonrest.CreateRequest[UserCreatePayload],
) (*commonrest.Response[commonrest.IdPayload], error) {
	id, err := u.userService.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
	}
	res := commonrest.IdPayload{ID: id}
	return commonrest.BuildResponse(res), nil
}

func (u *UserRoutes) update(
	ctx context.Context, req *commonrest.UpdateRequest[UserUpdatePayload],
) (*commonrest.EmptyResponse, error) {
	err := u.userService.Update(ctx, req.ID, req.Body.Name)
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(u.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	return &commonrest.EmptyResponse{}, nil
}
