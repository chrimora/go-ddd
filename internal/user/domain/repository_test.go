//go:build integration

package domain_test

import (
	"context"
	"goddd/internal/common/domain"
	commontest "goddd/internal/common/test"
	"goddd/internal/user/domain"
	"goddd/internal/user/test"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type UserSuite struct {
	suite.Suite
	app       *fx.App
	uf        test.UserFactory
	repo      domain.UserRepositoryI
	txManager *commondomain.TxManager
}

func (s *UserSuite) SetupSuite() {
	s.app = fx.New(
		test.IntegrationTestModule,
		fx.Populate(&s.uf),
		fx.Populate(&s.repo),
		fx.Populate(&s.txManager),
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
	ctx := commontest.TestContext()
	user := s.uf.Mock(s.T(), ctx, map[string]any{"Name": "Chris"})

	user, err := s.repo.Get(ctx, user.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Chris", user.Name())
}

func (s *UserSuite) TestRaceCondition() {
	ctx := commontest.TestContext()
	user := s.uf.Mock(s.T(), ctx)

	user.Update("Terry")
	err := s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, user)
	})
	require.NoError(s.T(), err)

	user.Update("Will")
	err = s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, user)
	})
	assert.ErrorIs(s.T(), err, commondomain.ErrRaceCondition)

	user, err = s.repo.Get(ctx, user.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Terry", user.Name())
}
