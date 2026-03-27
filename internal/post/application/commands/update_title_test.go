//go:build unit

package commands_test

import (
	"context"
	commontest "goddd/internal/common/test"
	"goddd/internal/post/application/commands"
	"goddd/internal/post/domain"
	"goddd/internal/post/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type PostUpdateTitleSuite struct {
	suite.Suite
	app     *fx.App
	pf      test.PostFactory
	repo    domain.PostRepositoryI
	command commands.UpdatePostTitleCommand
}

func (s *PostUpdateTitleSuite) SetupSuite() {
	s.app = fx.New(
		test.UnitTestModule,
		fx.Populate(&s.pf),
		fx.Populate(&s.repo),
		fx.Populate(&s.command),
	)
	s.app.Start(context.Background())
}
func (s *PostUpdateTitleSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestPostUpdateTitleSuite(t *testing.T) {
	suite.Run(t, new(PostUpdateTitleSuite))
}

func (s *PostUpdateTitleSuite) TestUpdateTitle() {
	ctx := commontest.TestContext()
	post := s.pf.Mock(s.T(), ctx, map[string]any{"Title": "Original Title"})

	_, err := s.command.Handle(ctx, commands.UpdatePostTitleInput{
		Id:    post.ID(),
		Title: "Updated Title",
	})
	require.NoError(s.T(), err)

	post, err = s.repo.Get(ctx, post.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Title", post.Title())
}
