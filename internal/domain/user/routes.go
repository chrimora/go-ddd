package user

import (
	"context"
	"errors"
	"gotemplate/internal/common"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UserServiceI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, UserCreatePayload) (uuid.UUID, error)
	Update(context.Context, uuid.UUID, UserUpdatePayload) error
}

type UserRoutes struct {
	log         *slog.Logger
	userService UserServiceI
}

func NewUserRoutes(log *slog.Logger, userService UserServiceI) *UserRoutes {
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
	common.IdPayload
	Name string `json:"name"`
}

func (u *UserRoutes) get(
	ctx context.Context, req *common.IdParam,
) (*common.Response[UserResponse], error) {
	user, err := u.userService.Get(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return nil, common.NotFoundResponse(u.log, ctx, err)
		default:
			return nil, common.UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	return common.BuildResponse(res), nil
}

func (u *UserRoutes) create(
	ctx context.Context, req *common.CreateRequest[UserCreatePayload],
) (*common.Response[common.IdPayload], error) {
	id, err := u.userService.Create(ctx, req.Body)
	if err != nil {
		return nil, common.UnexpectedErrorResponse(u.log, ctx, err)
	}
	res := common.IdPayload{ID: id}
	return common.BuildResponse(res), nil
}

func (u *UserRoutes) update(
	ctx context.Context, req *common.UpdateRequest[UserUpdatePayload],
) (*common.EmptyResponse, error) {
	err := u.userService.Update(ctx, req.ID, req.Body)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return nil, common.NotFoundResponse(u.log, ctx, err)
		default:
			return nil, common.UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	return &common.EmptyResponse{}, nil
}
