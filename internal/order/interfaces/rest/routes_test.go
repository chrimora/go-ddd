//go:build unit

package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"goddd/internal/common/infrastructure/middleware"
	commontest "goddd/internal/common/test"
	"goddd/internal/order/application/queries"
	"goddd/internal/order/domain"
	"goddd/internal/order/interfaces/rest"
	"goddd/internal/order/test"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

// GetOrdersByUserQuery requires ReadOrderSql (needs a real DB) - stub for unit tests
type noopGetOrdersByUser struct{}

func (s *noopGetOrdersByUser) Handle(_ context.Context, _ queries.GetOrdersByUserInput) (queries.GetOrdersByUserOutput, error) {
	return queries.GetOrdersByUserOutput{}, nil
}

type AddItemRouteSuite struct {
	suite.Suite
	app *fx.App
	of  test.OrderFactory
	api humatest.TestAPI
}

func (s *AddItemRouteSuite) SetupSuite() {
	var routes rest.OrderRoutes

	s.app = fx.New(
		test.UnitTestModule,
		fx.Supply(
			fx.Annotate(&noopGetOrdersByUser{}, fx.As(new(queries.GetOrdersByUserQuery))),
		),
		fx.Provide(
			rest.NewOrderRoutes,
		),
		fx.Populate(&s.of, &routes),
	)
	s.app.Start(context.Background())

	_, api := humatest.New(s.T())
	api.UseMiddleware(middleware.RequestIdMiddleware)
	api.UseMiddleware(middleware.UserAuthMiddleware)
	routes.Register(api)
	s.api = api
}

func (s *AddItemRouteSuite) TeardownSuite() {
	s.app.Stop(context.Background())
}

func TestAddItemRouteSuite(t *testing.T) {
	suite.Run(t, new(AddItemRouteSuite))
}

func (s *AddItemRouteSuite) TestAddItemSuccess() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)

	resp := s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", order.ID()), map[string]any{
		"name":     "Widget",
		"quantity": 1,
	})

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *AddItemRouteSuite) TestAddItemNotFound() {
	ctx := commontest.TestContext()

	resp := s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", uuid.New()), map[string]any{
		"name":     "Widget",
		"quantity": 1,
	})

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *AddItemRouteSuite) TestAddItemDuplicate() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)
	s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", order.ID()), map[string]any{"name": "Widget", "quantity": 1})

	resp := s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", order.ID()), map[string]any{
		"name":     "Widget",
		"quantity": 1,
	})

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
}

func (s *AddItemRouteSuite) TestAddItemOrderNotPending() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx, map[string]any{"Status": domain.Confirmed})

	resp := s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", order.ID()), map[string]any{
		"name":     "Widget",
		"quantity": 1,
	})

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
}

func (s *AddItemRouteSuite) TestAddItemValidation() {
	ctx := commontest.TestContext()
	order := s.of.Mock(s.T(), ctx)

	resp := s.api.PostCtx(ctx, fmt.Sprintf("/orders/%s/items", order.ID()), map[string]any{
		"name":     "Widget",
		"quantity": 0,
	})

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}
