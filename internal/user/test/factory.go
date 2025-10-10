package test

import (
	"context"
	common "goddd/internal/common/domain"
	"goddd/internal/common/test"
	"goddd/internal/user/application"
	"goddd/internal/user/domain"
	"testing"
)

type UserFactory struct {
	repo application.UserRepositoryI
}

func NewUserFactory(repo application.UserRepositoryI) *UserFactory {
	return &UserFactory{repo: repo}
}

func (f *UserFactory) Mock(t *testing.T, overrides ...map[string]any) *domain.User {
	ctx := context.Background()
	user := &domain.User{
		AggregateRoot: common.NewAggregateRoot(),
		Name:          "Christopher",
	}
	test.Merge(user, overrides)

	f.repo.Create(ctx, user)
	t.Cleanup(func() {
		f.repo.Remove(ctx, user)
	})
	return user
}
