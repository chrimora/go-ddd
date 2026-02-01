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
		log            *slog.Logger
		getUser        queries.GetUserQuery
		createUser     commands.CreateUserCommand
		userChangeName commands.UserChangeNameCommand
	}
)

func NewUserRoutes(
	log *slog.Logger,
	getUser queries.GetUserQuery,
	createUser commands.CreateUserCommand,
	userChangeName commands.UserChangeNameCommand,
) UserRoutes {
	return &userRoutes{
		log:            log,
		getUser:        getUser,
		createUser:     createUser,
		userChangeName: userChangeName,
	}
}

func (u *userRoutes) Register(api huma.API) {
	huma.Get(api, "/users/profiles", u.get)
	huma.Post(api, "/users/register", u.register)
	huma.Put(api, "/users/{id}/change-name", u.changeName)
}

type RegisterPayload struct {
	Name string `json:"name"`
}
type ChangeNamePayload struct {
	Name string `json:"name"`
}

type UserResponse struct {
	commonrest.IdPayload
	Name string `json:"name"`
}

func (u *userRoutes) get(
	ctx context.Context, req *commonrest.IdQuery,
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
	res.ID = user.ID()
	res.Name = user.Name()
	return commonrest.BuildResponse(res), nil
}

func (u *userRoutes) register(
	ctx context.Context, req *commonrest.CreateRequest[RegisterPayload],
) (*commonrest.Response[commonrest.IdPayload], error) {
	id, err := u.createUser.Handle(ctx, commands.CreateUserInput{Name: req.Body.Name})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
	}
	res := commonrest.IdPayload{ID: id}
	return commonrest.BuildResponse(res), nil
}

func (u *userRoutes) changeName(
	ctx context.Context, req *commonrest.UpdateRequest[ChangeNamePayload],
) (*commonrest.EmptyResponse, error) {
	_, err := u.userChangeName.Handle(ctx, commands.UserChangeNameInput{Id: req.ID, Name: req.Body.Name})
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
