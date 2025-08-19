package main

import (
	"context"
	"fmt"
	"gotemplate/internal/common"
	"gotemplate/internal/config"
	"gotemplate/internal/middleware"
	"gotemplate/internal/user"
	"net"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"go.uber.org/fx"
)

func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

func NewApi(routeCollection []common.RouteCollection, mux *http.ServeMux) *huma.API {
	api := humago.New(mux, huma.DefaultConfig("Go Template", "1.0"))
	api.UseMiddleware(middleware.RequestIdMiddleware)
	for _, routes := range routeCollection {
		routes.Register(api)
	}
	return &api
}

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			NewServeMux,
			fx.Annotate(
				NewApi,
				fx.ParamTags(`group:"routeCollection"`),
			),
		),
		config.Module,
		common.Module,
		user.Module,
		fx.Invoke(func(*http.Server, *huma.API) {}),
	).Run()
}
