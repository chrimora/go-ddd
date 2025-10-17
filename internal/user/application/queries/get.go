package queries

import (
	"context"
	commonapplication "goddd/internal/common/application"
	"goddd/internal/outbox"
	"goddd/internal/user/domain"
	"log/slog"

	"github.com/google/uuid"
)

type GetUserInput struct {
	Id uuid.UUID
}

type GetUserQuery commonapplication.QueryI[GetUserInput, *domain.User]

func NewGetUserQuery(
	log *slog.Logger,
	outboxRepo outbox.OutboxRepositoryI,
	userRepo domain.UserRepositoryI,
) GetUserQuery {
	return commonapplication.NewQuery(log, &getUser{
		userRepo: userRepo,
	})
}

type getUser struct {
	userRepo domain.UserRepositoryI
}

func (u *getUser) Handle(
	ctx context.Context, log *slog.Logger, input GetUserInput,
) (*domain.User, error) {
	user, err := u.userRepo.Get(ctx, input.Id)
	log.InfoContext(ctx, "Got", "user", user)
	return user, err
}
