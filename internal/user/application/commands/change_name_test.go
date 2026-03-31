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

type UserChangeNameSuite struct {
	suite.Suite
	app     *fx.App
	uf      test.UserFactory
	repo    domain.UserRepositoryI
	command commands.UserChangeNameCommand
}

func (s *UserChangeNameSuite) SetupSuite() {
	s.app = fx.New(
		test.UnitTestModule,
		fx.Populate(&s.uf),
		fx.Populate(&s.repo),
		fx.Populate(&s.command),
	)
	s.app.Start(context.Background())
}
func (s *UserChangeNameSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestUserChangeNameSuite(t *testing.T) {
	suite.Run(t, new(UserChangeNameSuite))
}

func (s *UserChangeNameSuite) TestChangeName() {
	ctx := commontest.TestContext()
	user := s.uf.Mock(s.T(), ctx, map[string]any{"Name": "Chris"})
	newName := "Christopher"

	_, err := s.command.Handle(ctx, commands.UserChangeNameInput{Id: user.ID(), Name: newName})
	require.NoError(s.T(), err)

	user, err = s.repo.Get(ctx, user.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), newName, user.Name())
}
