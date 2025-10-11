package commonrest

import (
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

type RouteCollection interface {
	Register(api huma.API)
}

func AsRouteCollection(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(RouteCollection)),
		fx.ResultTags(`group:"routeCollection"`),
	)
}
