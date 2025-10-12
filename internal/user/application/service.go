package application

import (
	"context"
	"goddd/internal/common/domain"
	"goddd/internal/outbox"
	"goddd/internal/user/domain"
	"log/slog"

	"github.com/google/uuid"
)

// TODO; split
// - Read: UserQueriesI
// - Write: UserCommandsI
type UserServiceI interface {
	Get(context.Context, uuid.UUID) (*domain.User, error)
	Create(context.Context, string) (uuid.UUID, error)
	Update(context.Context, uuid.UUID, string) error
}

type UserService UserServiceI
type userService struct {
	log        *slog.Logger
	txManager  *commondomain.TxManager
	outboxRepo outbox.OutboxRepositoryI
	userRepo   domain.UserRepositoryI
}

func NewUserService(
	log *slog.Logger,
	txManager *commondomain.TxManager,
	outboxRepo outbox.OutboxRepositoryI,
	userRepo domain.UserRepositoryI,
) UserService {
	return &userService{
		log:        log,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		userRepo:   userRepo,
	}
}

func (u *userService) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.Get(ctx, id)
	u.log.InfoContext(ctx, "Got", "user", user)
	return user, err
}

func (u *userService) Create(ctx context.Context, name string) (uuid.UUID, error) {
	user := domain.NewUser(name)
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

func (u *userService) Update(ctx context.Context, id uuid.UUID, name string) error {
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
