package test

import (
	"context"
	"goddd/internal/common/test"
	"goddd/internal/user/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

type UserFactory commontest.Mock[domain.User]
type userFactory struct {
	repo domain.UserRepositoryI
}

func NewUserFactory(repo domain.UserRepositoryI) UserFactory {
	return &userFactory{repo: repo}
}

func (f *userFactory) Mock(t *testing.T, overrides ...map[string]any) *domain.User {
	ctx := context.Background()
	fields := &struct {
		ID        uuid.UUID
		Version   int
		CreatedAt time.Time
		UpdatedAt time.Time
		Name      string
	}{
		ID:   uuid.New(),
		Name: "Christopher",
	}
	commontest.Merge(fields, overrides)

	user := domain.RehydrateUser(fields.ID, fields.Version, fields.CreatedAt, fields.UpdatedAt, fields.Name)
	f.repo.Create(ctx, user)
	t.Cleanup(func() {
		f.repo.Remove(ctx, user)
	})
	return user
}
