//go:build unit

package commands_test

import (
	"context"
	commontest "goddd/internal/common/test"
	"goddd/internal/user/application/commands"
	"goddd/internal/user/domain"
	"goddd/internal/user/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type UserCreateSuite struct {
	suite.Suite
	app     *fx.App
	repo    domain.UserRepositoryI
	command commands.CreateUserCommand
}

func (s *UserCreateSuite) SetupSuite() {
	s.app = fx.New(
		test.UnitTestModule,
		fx.Populate(&s.repo),
		fx.Populate(&s.command),
	)
	s.app.Start(context.Background())
}
func (s *UserCreateSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestUserCreateSuite(t *testing.T) {
	suite.Run(t, new(UserCreateSuite))
}

func (s *UserCreateSuite) TestCreate() {
	ctx := commontest.TestContext()
	name := "Christopher"

	id, err := s.command.Handle(ctx, commands.CreateUserInput{Name: name})
	require.NoError(s.T(), err)

	user, err := s.repo.Get(ctx, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), name, user.Name())
}
