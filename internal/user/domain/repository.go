package domain

import (
	"context"
	"database/sql"
	"errors"
	"goddd/internal/outbox"
	usersql "goddd/internal/user/infrastructure/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, pgx.Tx, *User) error
	Update(context.Context, pgx.Tx, *User) error
	Remove(context.Context, pgx.Tx, *User) error
}

type UserRepository UserRepositoryI

type userRepository struct {
	log        *slog.Logger
	userSql    *usersql.Queries
	outboxRepo outbox.OutboxRepositoryI
}

func NewUserRepository(
	log *slog.Logger,
	userSql *usersql.Queries,
	outboxRepo outbox.OutboxRepositoryI,
) UserRepository {
	return &userRepository{
		log:        log,
		userSql:    userSql,
		outboxRepo: outboxRepo,
	}
}

func (u *userRepository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := u.userSql.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound(id)
		default:
			return nil, err
		}
	}
	return RehydrateUser(
		user.ID, int(user.Version), user.CreatedAt, user.UpdatedAt, user.Name,
	), nil
}

func (u *userRepository) Create(ctx context.Context, tx pgx.Tx, user *User) error {
	_, err := u.userSql.WithTx(tx).CreateUser(
		ctx,
		usersql.CreateUserParams{
			ID:        user.ID(),
			Version:   int32(user.Version()),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
			Name:      user.Name(),
		},
	)
	if err != nil {
		return err
	}
	return u.outboxRepo.CreateMany(ctx, tx, user.PullEvents()...)
}

func (u *userRepository) Update(ctx context.Context, tx pgx.Tx, user *User) error {
	_, err := u.userSql.WithTx(tx).UpdateUser(
		ctx,
		usersql.UpdateUserParams{
			ID:        user.ID(),
			Version:   int32(user.Version()),
			UpdatedAt: user.UpdatedAt(),
			Name:      user.Name(),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRaceCondition(user.ID())
		default:
			return err
		}
	}
	return u.outboxRepo.CreateMany(ctx, tx, user.PullEvents()...)
}

func (u *userRepository) Remove(ctx context.Context, tx pgx.Tx, user *User) error {
	return u.userSql.WithTx(tx).RemoveUser(ctx, user.ID())
}
