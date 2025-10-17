package commands

import (
	"context"
	commonapplication "goddd/internal/common/application"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/outbox"
	"goddd/internal/user/domain"
	"log/slog"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	Name string
}

type CreateUserCommand commonapplication.CommandI[CreateUserInput]

func NewCreateUserCommand(
	log *slog.Logger,
	txManager *commondomain.TxManager,
	outboxRepo outbox.OutboxRepositoryI,
	userRepo domain.UserRepositoryI,
) CreateUserCommand {
	return commonapplication.NewCommand(log, &createUser{
		txManager:  txManager,
		outboxRepo: outboxRepo,
		userRepo:   userRepo,
	})
}

type createUser struct {
	txManager  *commondomain.TxManager
	outboxRepo outbox.OutboxRepositoryI
	userRepo   domain.UserRepositoryI
}

func (u *createUser) Handle(
	ctx context.Context, log *slog.Logger, input CreateUserInput,
) (uuid.UUID, error) {
	user := domain.NewUser(input.Name)
	log.InfoContext(ctx, "Creating", "user", user)

	err := u.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		err := u.userRepo.Create(txCtx, user)
		if err != nil {
			return err
		}
		return u.outboxRepo.CreateMany(txCtx, user.PullEvents()...)
	})
	return user.ID(), err
}
