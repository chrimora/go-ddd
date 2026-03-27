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
		common.APIModule,
		config.APIModule,
		outbox.APIModule,
		post.APIModule,
		fx.Invoke(func(*http.Server, *huma.API) {}),
	).Run()
}
