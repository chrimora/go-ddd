package commands

import (
	"context"
	commonapplication "goddd/internal/common/application"
	commondomain "goddd/internal/common/domain"
	"goddd/internal/user/domain"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UpdateUserInput struct {
	Id   uuid.UUID
	Name string
}

type UpdateUserCommand commonapplication.CommandI[UpdateUserInput]

func NewUpdateUserCommand(
	log *slog.Logger,
	txManager *commondomain.TxManager,
	userRepo domain.UserRepositoryI,
) UpdateUserCommand {
	return commonapplication.NewCommand(log, &updateUser{
		txManager: txManager,
		userRepo:  userRepo,
	})
}

type updateUser struct {
	txManager *commondomain.TxManager
	userRepo  domain.UserRepositoryI
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

	err = u.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return u.userRepo.Update(ctx, tx, user)
	})
	return user.ID(), err
}
