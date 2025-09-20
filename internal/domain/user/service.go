package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type UserRepositoryI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, *User) error
	Update(context.Context, *User) error
}

type UserService struct {
	log      *slog.Logger
	userRepo UserRepositoryI
}

func NewUserService(log *slog.Logger, userRepo UserRepositoryI) UserServiceI {
	return &UserService{
		log:      log,
		userRepo: userRepo,
	}
}

func (u *UserService) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserService) Create(ctx context.Context, req UserCreatePayload) (uuid.UUID, error) {
	user := NewUser(req.Name)

	u.log.InfoContext(ctx, "Creating", "user", user)
	err := u.userRepo.Create(ctx, user)
	return user.ID, err
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
