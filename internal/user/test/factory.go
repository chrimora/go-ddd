package test

import (
	"context"
	"goddd/internal/common/domain"
	"goddd/internal/common/test"
	"goddd/internal/user/domain"
	"testing"
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
	user := &domain.User{
		AggregateRoot: commondomain.NewAggregateRoot(),
		Name:          "Christopher",
	}
	commontest.Merge(user, overrides)

	f.repo.Create(ctx, user)
	t.Cleanup(func() {
		f.repo.Remove(ctx, user)
	})
	return user
}
