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

type CreateUserInput struct {
	Name string
}

type CreateUserCommand commonapplication.CommandI[CreateUserInput]

func NewCreateUserCommand(
	log *slog.Logger,
	txManager commondomain.TxManager,
	userRepo domain.UserRepositoryI,
) CreateUserCommand {
	return commonapplication.NewCommand(log, &createUser{
		txManager: txManager,
		userRepo:  userRepo,
	})
}

type createUser struct {
	txManager commondomain.TxManager
	userRepo  domain.UserRepositoryI
}

func (u *createUser) Handle(
	ctx context.Context, log *slog.Logger, input CreateUserInput,
) (uuid.UUID, error) {
	user := domain.NewUser(input.Name)
	log.InfoContext(ctx, "Creating", "user", user)

	err := u.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return u.userRepo.Create(ctx, tx, user)
	})
	return user.ID(), err
}
