//go:build integration

package commands_test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	commontest "goddd/internal/common/test"
	"goddd/internal/order/application/commands"
	"goddd/internal/order/domain"
	"goddd/internal/order/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type AddItemIntegrationSuite struct {
	suite.Suite
	app     *fx.App
	of      test.OrderFactory
	repo    domain.OrderRepositoryI
	command commands.AddOrderItemCommand
}

func (s *AddItemIntegrationSuite) SetupSuite() {
	s.app = fx.New(
		test.IntegrationTestModule,
		fx.Populate(&s.of),
		fx.Populate(&s.repo),
		fx.Populate(&s.command),
	)
	s.app.Start(context.Background())
}

func (s *AddItemIntegrationSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestAddItemIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AddItemIntegrationSuite))
}

func (s *AddItemIntegrationSuite) TestNotFound() {
	ctx := commontest.TestContext()

	_, err := s.command.Handle(ctx, commands.AddOrderItemInput{
		OrderId:  commondomain.NewUUID(),
		Name:     "Widget",
		Quantity: 1,
	})

	assert.ErrorIs(s.T(), err, commondomain.ErrNotFound)
}

func (s *AddItemIntegrationSuite) TestOrderNotPendingDoesNotUpdate() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx, map[string]any{"Status": domain.Confirmed})

	_, err := s.command.Handle(ctx, commands.AddOrderItemInput{
		OrderId:  order.ID(),
		Name:     "Widget",
		Quantity: 1,
	})
	assert.ErrorIs(s.T(), err, domain.ErrOrderNotPending)

	result, err := s.repo.Get(ctx, order.ID())
	require.NoError(s.T(), err)
	assert.Empty(s.T(), result.Items())
}

func (s *AddItemIntegrationSuite) TestAddItemPersists() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)

	_, err := s.command.Handle(ctx, commands.AddOrderItemInput{
		OrderId:  order.ID(),
		Name:     "Widget",
		Quantity: 3,
	})
	require.NoError(s.T(), err)

	_, err = s.repo.Get(ctx, order.ID())
	require.NoError(s.T(), err)
}
