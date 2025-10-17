package rest

import (
	"context"
	"errors"
	"goddd/internal/common/domain"
	"goddd/internal/common/interfaces/rest"
	"goddd/internal/user/application/commands"
	"goddd/internal/user/application/queries"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

type (
	UserRoutes commonrest.RouteCollection
	userRoutes struct {
		log        *slog.Logger
		getUser    queries.GetUserQuery
		createUser commands.CreateUserCommand
		updateUser commands.UpdateUserCommand
	}
)

func NewUserRoutes(
	log *slog.Logger,
	getUser queries.GetUserQuery,
	createUser commands.CreateUserCommand,
	updateUser commands.UpdateUserCommand,
) UserRoutes {
	return &userRoutes{
		log:        log,
		getUser:    getUser,
		createUser: createUser,
		updateUser: updateUser,
	}
}

func (u *userRoutes) Register(api huma.API) {
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

func (u *userRoutes) get(
	ctx context.Context, req *commonrest.IdParam,
) (*commonrest.Response[UserResponse], error) {
	user, err := u.getUser.Handle(ctx, queries.GetUserInput{Id: req.ID})
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

func (u *userRoutes) create(
	ctx context.Context, req *commonrest.CreateRequest[UserCreatePayload],
) (*commonrest.Response[commonrest.IdPayload], error) {
	id, err := u.createUser.Handle(ctx, commands.CreateUserInput{Name: req.Body.Name})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
	}
	res := commonrest.IdPayload{ID: id}
	return commonrest.BuildResponse(res), nil
}

func (u *userRoutes) update(
	ctx context.Context, req *commonrest.UpdateRequest[UserUpdatePayload],
) (*commonrest.EmptyResponse, error) {
	_, err := u.updateUser.Handle(ctx, commands.UpdateUserInput{Id: req.ID, Name: req.Body.Name})
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
