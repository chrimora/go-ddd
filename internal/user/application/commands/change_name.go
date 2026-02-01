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

type UserChangeNameInput struct {
	Id   uuid.UUID
	Name string
}

type UserChangeNameCommand commonapplication.CommandI[UserChangeNameInput]

func NewUserChangeNameCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	userRepo domain.UserRepositoryI,
) UserChangeNameCommand {
	return commonapplication.NewCommand(log, &changeName{
		txManager: txManager,
		userRepo:  userRepo,
	})
}

type changeName struct {
	txManager commondomain.TxManager
	userRepo  domain.UserRepositoryI
}

func (u *changeName) Handle(
	ctx context.Context, log *slog.Logger, input UserChangeNameInput,
) (uuid.UUID, error) {
	user, err := u.userRepo.Get(ctx, input.Id)
	if err != nil {
		return user.ID(), err
	}

	user.ChangeName(input.Name)
	log.InfoContext(ctx, "Updating", "user", user)

	err = u.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return u.userRepo.Update(ctx, tx, user)
	})
	return user.ID(), err
}
