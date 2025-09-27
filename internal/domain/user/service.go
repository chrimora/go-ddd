package user

import (
	"context"
	"goddd/internal/domain/common"
	"goddd/internal/domain/outbox"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepositoryI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
	Remove(context.Context, *User) error
}

// TODO; split
// - Read: UserQueriesI
// - Write: UserCommandsI
type UserServiceI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, string) (uuid.UUID, error)
	Update(context.Context, uuid.UUID, string) error
}

type UserService struct {
	log        *slog.Logger
	txManager  *common.TxManager
	outboxRepo outbox.OutboxRepositoryI
	userRepo   UserRepositoryI
}

func NewUserService(
	log *slog.Logger,
	txManager *common.TxManager,
	outboxRepo outbox.OutboxRepositoryI,
	userRepo UserRepositoryI,
) *UserService {
	return &UserService{
		log:        log,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		userRepo:   userRepo,
	}
}

func (u *UserService) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := u.userRepo.Get(ctx, id)
	u.log.InfoContext(ctx, "Got", "user", user)
	return user, err
}

func (u *UserService) Create(ctx context.Context, name string) (uuid.UUID, error) {
	user := NewUser(name)
	u.log.InfoContext(ctx, "Creating", "user", user)

	err := u.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		err := u.userRepo.Create(txCtx, user)
		if err != nil {
			return err
		}
		return u.outboxRepo.CreateMany(txCtx, user.PullEvents()...)
	})

	return user.ID, err
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, name string) error {
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	user.Update(name)
	u.log.InfoContext(ctx, "Updating", "user", user)

	return u.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		err := u.userRepo.Update(txCtx, user)
		if err != nil {
			return err
		}
		return u.outboxRepo.CreateMany(txCtx, user.PullEvents()...)
	})
}
