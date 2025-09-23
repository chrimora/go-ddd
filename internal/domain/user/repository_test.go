//go:build integration

package user_test

import (
	"context"
	"goddd/internal/domain"
	"goddd/internal/domain/common"
	"goddd/test"
	"goddd/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type UserSuite struct {
	suite.Suite
	app  *fx.App
	uf   *factories.UserFactory
	repo domain.UserRepositoryI
}

func (s *UserSuite) SetupSuite() {
	s.app = fx.New(
		test.IntegrationTestModule,
		fx.Populate(&s.uf),
		fx.Populate(&s.repo),
	)
	s.app.Start(context.Background())
}
func (s *UserSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) TestGet() {
	user := s.uf.Mock(s.T(), map[string]any{"Name": "Chris"})

	user, err := s.repo.Get(context.Background(), user.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Chris", user.Name)
}

func (s *UserSuite) TestRaceCondition() {
	ctx := context.Background()
	user := s.uf.Mock(s.T())

	user.Update("Terry")
	err := s.repo.Update(ctx, user)
	require.NoError(s.T(), err)

	user.Update("Will")
	err = s.repo.Update(ctx, user)
	assert.ErrorIs(s.T(), err, common.ErrRaceCondition)

	user, err = s.repo.Get(ctx, user.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Terry", user.Name)
}
