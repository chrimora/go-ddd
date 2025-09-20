package user

import (
	"context"
	"database/sql"
	"errors"
	"gotemplate/internal/domain/common"
	userdb "gotemplate/internal/domain/user/db"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepository struct {
	log       *slog.Logger
	txFactory common.TxFactory
	userSql   *userdb.Queries
}

func NewUserRepository(log *slog.Logger, txFactory common.TxFactory, userSql *userdb.Queries) UserRepositoryI {
	return &UserRepository{
		log:       log,
		txFactory: txFactory,
		userSql:   userSql,
	}
}

func fromDB(user userdb.User) *User {
	return &User{
		AggregateRoot: common.AggregateRoot{ID: user.ID},
		UpdatedAt:     user.UpdatedAt,
		Name:          user.Name,
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

func (u *UserRepository) Create(ctx context.Context, user *User) error {
	tx, err := u.txFactory(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = u.userSql.WithTx(tx.Tx).CreateUser(ctx, userdb.CreateUserParams{ID: user.ID, Name: user.Name})

	tx.TrackEvents(user)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (u *UserRepository) Update(ctx context.Context, user *User) error {
	_, err := u.userSql.UpdateUser(
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
