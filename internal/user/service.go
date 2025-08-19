package user

import (
	"context"
	"database/sql"
	"errors"
	userdb "gotemplate/internal/user/db"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepositoryI interface {
	GetUser(context.Context, uuid.UUID) (userdb.User, error)
	CreateUser(context.Context, string) (uuid.UUID, error)
	UpdateUser(context.Context, userdb.UpdateUserParams) error
}

type UserService struct {
	log      *slog.Logger
	userRepo UserRepositoryI
}

func NewUserService(log *slog.Logger, userRepo UserRepositoryI) UserServiceI {
	return &UserService{
		log:      log,
		userRepo: userRepo,
	}
}

func (u *UserService) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound(id)
		default:
			return nil, err
		}
	}
	return &User{ID: user.ID, Name: user.Name}, nil
}

func (u *UserService) Create(ctx context.Context, req UserCreatePayload) (uuid.UUID, error) {
	id, err := u.userRepo.CreateUser(ctx, req.Name)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, req UserUpdatePayload) error {
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound(id)
		default:
			return err
		}
	}
	err = u.userRepo.UpdateUser(
		ctx,
		userdb.UpdateUserParams{ID: id, UpdatedAt: user.UpdatedAt, Name: req.Name},
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRaceCondition(id)
		default:
			return err
		}
	}
	return err
}
