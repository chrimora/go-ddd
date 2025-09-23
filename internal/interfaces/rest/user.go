package rest

import (
	"context"
	"errors"
	"goddd/internal/domain"
	"goddd/internal/domain/common"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

type UserRoutes struct {
	log         *slog.Logger
	userService domain.UserServiceI
}

func NewUserRoutes(log *slog.Logger, userService domain.UserServiceI) *UserRoutes {
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
	IdPayload
	Name string `json:"name"`
}

func (u *UserRoutes) get(
	ctx context.Context, req *IdParam,
) (*Response[UserResponse], error) {
	user, err := u.userService.Get(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return nil, NotFoundResponse(u.log, ctx, err)
		default:
			return nil, UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	return BuildResponse(res), nil
}

func (u *UserRoutes) create(
	ctx context.Context, req *CreateRequest[UserCreatePayload],
) (*Response[IdPayload], error) {
	id, err := u.userService.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, UnexpectedErrorResponse(u.log, ctx, err)
	}
	res := IdPayload{ID: id}
	return BuildResponse(res), nil
}

func (u *UserRoutes) update(
	ctx context.Context, req *UpdateRequest[UserUpdatePayload],
) (*EmptyResponse, error) {
	err := u.userService.Update(ctx, req.ID, req.Body.Name)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return nil, NotFoundResponse(u.log, ctx, err)
		default:
			return nil, UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	return &EmptyResponse{}, nil
}
