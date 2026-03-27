package commonrest

import (
	"context"
	"goddd/internal/common/infrastructure/middleware"
	"goddd/internal/config"
	"log/slog"
	"net"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"go.uber.org/fx"
)

func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

func NewApi(routeCollection []RouteCollection, mux *http.ServeMux) *huma.API {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	api := humago.New(mux, huma.DefaultConfig("Go Template", "1.0"))
	api.UseMiddleware(middleware.RequestIdMiddleware)
	api.UseMiddleware(middleware.UserAuthMiddleware)
	for _, routes := range routeCollection {
		routes.Register(api)
	}
	return &api
}

func NewHTTPServer(lc fx.Lifecycle, cfg *config.ServerConfig, mux *http.ServeMux, log *slog.Logger) *http.Server {
	srv := &http.Server{Addr: cfg.Port, Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Server start", "address", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Server stopping")
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
