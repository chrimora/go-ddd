//go:build integration

package domain_test

import (
	"context"
	commondomain "goddd/internal/common/domain"
	commontest "goddd/internal/common/test"
	"goddd/internal/order/domain"
	"goddd/internal/order/test"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type OrderSuite struct {
	suite.Suite
	app       *fx.App
	of        test.OrderFactory
	repo      domain.OrderRepositoryI
	txManager commondomain.TxManager
}

func (s *OrderSuite) SetupSuite() {
	s.app = fx.New(
		test.IntegrationTestModule,
		fx.Populate(&s.of),
		fx.Populate(&s.repo),
		fx.Populate(&s.txManager),
	)
	s.app.Start(context.Background())
}

func (s *OrderSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (s *OrderSuite) TestGetNotFound() {
	ctx := commontest.TestContext()

	_, err := s.repo.Get(ctx, commondomain.NewUUID())

	assert.ErrorIs(s.T(), err, commondomain.ErrNotFound)
}

func (s *OrderSuite) TestGet() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx, map[string]any{"Status": domain.Confirmed})

	result, err := s.repo.Get(ctx, order.ID())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), domain.Confirmed, result.Status())
}

func (s *OrderSuite) TestRaceCondition() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)

	err := s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, order)
	})
	require.NoError(s.T(), err)

	err = s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, order)
	})
	assert.ErrorIs(s.T(), err, commondomain.ErrRaceCondition)
}

func (s *OrderSuite) TestUpdatePersistsItems() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)

	_ = order.AddItem("Widget", 2, 500)
	_ = order.AddItem("Gadget", 1, 1000)
	err := s.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		return s.repo.Update(ctx, tx, order)
	})
	require.NoError(s.T(), err)

	result, err := s.repo.Get(ctx, order.ID())
	require.NoError(s.T(), err)
	assert.Len(s.T(), result.Items(), 2)
	item1 := result.Items()[0]
	item2 := result.Items()[1]
	assert.Equal(s.T(), "Widget", item1.Name())
	assert.Equal(s.T(), 2, item1.Quantity())
	assert.Equal(s.T(), int64(500), item1.UnitPrice())
	assert.Equal(s.T(), "Gadget", item2.Name())
	assert.Equal(s.T(), 1, item2.Quantity())
	assert.Equal(s.T(), int64(1000), item2.UnitPrice())
}
