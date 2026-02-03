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
	"github.com/google/uuid"
)

type (
	UserRoutes commonrest.RouteCollection
	userRoutes struct {
		log            *slog.Logger
		getUser        queries.GetUserQuery
		getProfiles    queries.GetProfilesQuery
		createUser     commands.CreateUserCommand
		userChangeName commands.UserChangeNameCommand
	}
)

func NewUserRoutes(
	log *slog.Logger,
	getUser queries.GetUserQuery,
	getProfiles queries.GetProfilesQuery,
	createUser commands.CreateUserCommand,
	userChangeName commands.UserChangeNameCommand,
) UserRoutes {
	return &userRoutes{
		log:            log,
		getUser:        getUser,
		getProfiles:    getProfiles,
		createUser:     createUser,
		userChangeName: userChangeName,
	}
}

func (u *userRoutes) Register(api huma.API) {
	huma.Get(api, "/users/{id}/profile", u.get)
	huma.Get(api, "/users/profiles", u.get_profiles)
	huma.Post(api, "/users/register", u.register)
	huma.Put(api, "/users/{id}/change-name", u.changeName)
}

type UserProfile struct {
	commonrest.IdPayload
	Name string `json:"name"`
}

func (u *userRoutes) get(
	ctx context.Context, req *commonrest.IdParam,
) (*commonrest.Response[UserProfile], error) {
	user, err := u.getUser.Handle(ctx, queries.GetUserInput{Id: req.ID})
	if err != nil {
		switch {
		case errors.Is(err, commondomain.ErrNotFound):
			return nil, commonrest.NotFoundResponse(u.log, ctx, err)
		default:
			return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
		}
	}
	res := UserProfile{}
	res.ID = user.ID()
	res.Name = user.Name()
	return commonrest.BuildResponse(res), nil
}

type ProfilesQuery struct {
	commonrest.PaginationQuery
	Name string `query:"name"`
}

func (u *userRoutes) get_profiles(
	ctx context.Context, req *ProfilesQuery,
) (*commonrest.Response[commonrest.Page[UserProfile]], error) {
	var after *uuid.UUID
	if req.AfterCursor != uuid.Nil {
		after = &req.AfterCursor
	}
	var name *string
	if req.Name != "" {
		name = &req.Name
	}

	out, err := u.getProfiles.Handle(ctx, queries.GetProfilesInput{
		Limit: req.Limit,
		After: after,
		Name:  name,
	})
	if err != nil {
		return nil, commonrest.UnexpectedErrorResponse(u.log, ctx, err)
	}

	res := commonrest.Page[UserProfile]{}
	items := make([]UserProfile, len(out.Profiles))
	for i, user := range out.Profiles {
		profile := UserProfile{}
		profile.ID = user.ID()
		profile.Name = user.Name()
		items[i] = profile
	}
	res.Items = items
	res.NextCursor = out.Next
	return commonrest.BuildResponse(res), nil
}

type RegisterPayload struct {
	Name string `json:"name"`
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

type ChangeNamePayload struct {
	Name string `json:"name"`
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
