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

// TODO; split
// - Read: UserQueriesI
// - Write: UserCommandsI
type UserServiceI interface {
	Get(context.Context, uuid.UUID) (*User, error)
	Create(context.Context, string) (uuid.UUID, error)
	Update(context.Context, uuid.UUID, string) error
}

type UserService struct {
	log      *slog.Logger
	userRepo UserRepositoryI
}

func NewUserService(log *slog.Logger, userRepo UserRepositoryI) *UserService {
	return &UserService{
		log:      log,
		userRepo: userRepo,
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
	err := u.userRepo.Create(ctx, user)
	return user.ID, err
}

func (u *UserService) Update(ctx context.Context, id uuid.UUID, name string) error {
	user, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	u.log.InfoContext(ctx, "Updating", "user", user)
	user.Update(name)
	return u.userRepo.Update(ctx, user)
}
