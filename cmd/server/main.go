package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/post"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		common.ServerModule,
		config.ServerModule,
		outbox.ServerModule,
		post.ServerModule,
		fx.Invoke(func(*http.Server, *huma.API) {}),
	).Run()
}
