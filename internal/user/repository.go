package user

import (
	"context"
	"database/sql"
	"errors"
	userdb "gotemplate/internal/user/db"
	"log/slog"

	"github.com/google/uuid"
)

type UserSqlI interface {
	GetUser(context.Context, uuid.UUID) (userdb.User, error)
	CreateUser(context.Context, string) (uuid.UUID, error)
	UpdateUser(context.Context, userdb.UpdateUserParams) error
}

type UserRepository struct {
	log     *slog.Logger
	userSql UserSqlI
	// eventsSql
}

func NewUserRepository(log *slog.Logger, userSql UserSqlI) UserRepositoryI {
	return &UserRepository{
		log:     log,
		userSql: userSql,
		// eventsSql: eventsSql,
	}
}

func (u *UserRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := u.userSql.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound(id)
		default:
			return nil, err
		}
	}
	return fromDB(user), nil
}

func (u *UserRepository) Create(ctx context.Context, user *User) (uuid.UUID, error) {
	id, err := u.userSql.CreateUser(ctx, user.Name)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (u *UserRepository) Update(ctx context.Context, user *User) error {
	err := u.userSql.UpdateUser(
		ctx,
		userdb.UpdateUserParams{ID: user.ID, UpdatedAt: user.UpdatedAt, Name: user.Name},
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRaceCondition(user.ID)
		default:
			return err
		}
	}
	return err
}
