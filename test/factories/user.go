package factories

import (
	"context"
	"gotemplate/internal/domain"
	"gotemplate/internal/domain/common"
	"testing"
)

type UserFactory struct {
	repo domain.UserRepositoryI
}

func NewUserFactory(repo domain.UserRepositoryI) *UserFactory {
	return &UserFactory{repo: repo}
}

func (f *UserFactory) Mock(t *testing.T, overrides ...map[string]any) *domain.User {
	ctx := context.Background()
	user := &domain.User{
		AggregateRoot: common.NewAggregateRoot(),
		Name:          "Christopher",
	}
	merge(user, overrides)

	f.repo.Create(ctx, user)
	t.Cleanup(func() {
		f.repo.Remove(ctx, user)
	})
	return user
}
