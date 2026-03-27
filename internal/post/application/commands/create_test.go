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

type PostCreateSuite struct {
	suite.Suite
	app     *fx.App
	repo    domain.PostRepositoryI
	command commands.CreatePostCommand
}

func (s *PostCreateSuite) SetupSuite() {
	s.app = fx.New(
		test.UnitTestModule,
		fx.Populate(&s.repo),
		fx.Populate(&s.command),
	)
	s.app.Start(context.Background())
}
func (s *PostCreateSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestPostCreateSuite(t *testing.T) {
	suite.Run(t, new(PostCreateSuite))
}

func (s *PostCreateSuite) TestCreate() {
	ctx := commontest.TestContext()

	id, err := s.command.Handle(ctx, commands.CreatePostInput{Title: "Hello World", Author: "alice"})
	require.NoError(s.T(), err)

	post, err := s.repo.Get(ctx, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Hello World", post.Title())
	assert.Equal(s.T(), "alice", post.Author())
}
