package user

import (
	"context"
	"database/sql"
	"errors"
	"gotemplate/internal/domain/common"
	userdb "gotemplate/internal/infrastructure/sql/codegen/user"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepository struct {
	log       *slog.Logger
	txFactory common.TxFactory
	userSql   *userdb.Queries
}

func NewUserRepository(log *slog.Logger, txFactory common.TxFactory, userSql *userdb.Queries) *UserRepository {
	return &UserRepository{
		log:       log,
		txFactory: txFactory,
		userSql:   userSql,
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
		AggregateRoot: common.NewAggregateRootFromFields(
			user.ID, int(user.Version), user.CreatedAt, user.UpdatedAt,
		),
		Name: user.Name,
	}, nil
}

func (u *UserRepository) Create(ctx context.Context, user *User) error {
	tx, err := u.txFactory(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = u.userSql.WithTx(tx.Tx).CreateUser(
		ctx,
		userdb.CreateUserParams{
			ID:        user.ID,
			Version:   int32(user.GetVersion()),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		},
	)

	tx.AddEvents(user.PullEvents()...)
	return tx.Commit()
}

func (u *UserRepository) Update(ctx context.Context, user *User) error {
	_, err := u.userSql.UpdateUser(
		ctx,
		userdb.UpdateUserParams{
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
	return u.userSql.RemoveUser(ctx, user.ID)
}
