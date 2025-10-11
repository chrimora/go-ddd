package domain

import (
	"context"
	"database/sql"
	"errors"
	"goddd/internal/common/domain"
	usersql "goddd/internal/user/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepository struct {
	log     *slog.Logger
	userSql *usersql.Queries
}

func NewUserRepository(
	log *slog.Logger,
	userSql *usersql.Queries,
) *UserRepository {
	return &UserRepository{
		log:     log,
		userSql: userSql,
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
	return &User{
		AggregateRoot: commondomain.NewAggregateRootFromFields(
			user.ID, int(user.Version), user.CreatedAt, user.UpdatedAt,
		),
		Name: user.Name,
	}, nil
}
func (u *UserRepository) Create(ctx context.Context, user *User) error {
	_, err := commondomain.WithTxFromCtx(u.userSql, ctx).CreateUser(
		ctx,
		usersql.CreateUserParams{
			ID:        user.ID,
			Version:   int32(user.GetVersion()),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		},
	)
	return err
}

func (u *UserRepository) Update(ctx context.Context, user *User) error {
	_, err := commondomain.WithTxFromCtx(u.userSql, ctx).UpdateUser(
		ctx,
		usersql.UpdateUserParams{
			ID:        user.ID,
			Version:   int32(user.GetVersion()),
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		},
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

func (u *UserRepository) Remove(ctx context.Context, user *User) error {
	return commondomain.WithTxFromCtx(u.userSql, ctx).RemoveUser(ctx, user.ID)
}
