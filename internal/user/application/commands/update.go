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

type UpdateUserInput struct {
	Id   uuid.UUID
	Name string
}

type UpdateUserCommand commonapplication.CommandI[UpdateUserInput]

func NewUpdateUserCommand(
	log *slog.Logger,
	txManager *commondomain.TxManager,
	outboxRepo outbox.OutboxRepositoryI,
	userRepo domain.UserRepositoryI,
) UpdateUserCommand {
	return commonapplication.NewCommand(log, &updateUser{
		txManager:  txManager,
		outboxRepo: outboxRepo,
		userRepo:   userRepo,
	})
}

type updateUser struct {
	txManager  *commondomain.TxManager
	outboxRepo outbox.OutboxRepositoryI
	userRepo   domain.UserRepositoryI
}

func (u *updateUser) Handle(
	ctx context.Context, log *slog.Logger, input UpdateUserInput,
) (uuid.UUID, error) {
	user, err := u.userRepo.Get(ctx, input.Id)
	if err != nil {
		return user.ID(), err
	}

	user.Update(input.Name)
	log.InfoContext(ctx, "Updating", "user", user)

	err = u.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		err := u.userRepo.Update(txCtx, user)
		if err != nil {
			return err
		}
		return u.outboxRepo.CreateMany(txCtx, user.PullEvents()...)
	})
	return user.ID(), err
}
