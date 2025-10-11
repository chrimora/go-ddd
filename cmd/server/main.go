package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/user"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "server"}
	fx.New(
		fx.Supply(service),
		common.ServerModule,
		config.Module,
		outbox.ServerModule,
		user.ServerModule,
		fx.Invoke(func(*http.Server, *huma.API) {}),
	).Run()
}
