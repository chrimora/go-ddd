package user

import (
	"context"
	"gotemplate/internal/domain/common"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepositoryI interface {
	WithTx(*common.Tx) UserRepositoryI
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
}

type UserService struct {
	log       *slog.Logger
	txFactory common.TxFactory
	userRepo  UserRepositoryI
}

func NewUserService(log *slog.Logger, txFactory common.TxFactory, userRepo UserRepositoryI) UserServiceI {
	return &UserService{
		log:       log,
		txFactory: txFactory,
		userRepo:  userRepo,
	}
}

func (u *UserService) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserService) Create(ctx context.Context, req UserCreatePayload) (uuid.UUID, error) {
	user := NewUser(req.Name)

	tx, err := u.txFactory(ctx)
	if err != nil {
		return user.ID, err
	}
	defer tx.Rollback()

	err = u.userRepo.WithTx(tx).Create(ctx, user)
	if err != nil {
		return user.ID, err
	}

	tx.TrackEvents(user)
	err = tx.Commit()
	return user.ID, nil
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, req UserUpdatePayload) error {
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	user.Update(req.Name)

	u.log.InfoContext(ctx, "Updating", "user", user)
	return u.userRepo.Update(ctx, user)
}
