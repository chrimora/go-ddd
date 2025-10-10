package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/domain"
	"goddd/internal/infrastructure/sql"
	"goddd/internal/interfaces/rest"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "server"}
	fx.New(
		fx.Supply(service),
		fx.Provide(
			common.NewLogger,
			rest.NewHTTPServer,
			rest.NewServeMux,
			fx.Annotate(
				rest.NewApi,
				fx.ParamTags(`group:"routeCollection"`),
			),
		),
		config.Module,
		domain.Module,
		sql.Module,
		rest.Module,
		fx.Invoke(func(*http.Server, *huma.API) {}),
	).Run()
}
