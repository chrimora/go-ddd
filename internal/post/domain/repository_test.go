//go:build integration

package domain_test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	commontest "goddd/internal/common/test"
	"goddd/internal/post/domain"
	"goddd/internal/post/test"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type PostSuite struct {
	suite.Suite
	app       *fx.App
	pf        test.PostFactory
	repo      domain.PostRepositoryI
	txManager commondomain.TxManager
}

func (s *PostSuite) SetupSuite() {
	s.app = fx.New(
		test.IntegrationTestModule,
		fx.Populate(&s.pf),
		fx.Populate(&s.repo),
		fx.Populate(&s.txManager),
	)
	s.app.Start(context.Background())
}
func (s *PostSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestPostSuite(t *testing.T) {
	suite.Run(t, new(PostSuite))
}

func (s *PostSuite) TestGet() {
	ctx := commontest.TestContext()
	post := s.pf.Mock(s.T(), ctx, map[string]any{"Title": "Hello World"})

	post, err := s.repo.Get(ctx, post.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Hello World", post.Title())
}

func (s *PostSuite) TestRaceCondition() {
	ctx := commontest.TestContext()
	post := s.pf.Mock(s.T(), ctx)

	err := s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, post)
	})
	require.NoError(s.T(), err)

	err = s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, post)
	})
	assert.ErrorIs(s.T(), err, commondomain.ErrRaceCondition)
}
