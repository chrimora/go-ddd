package middleware

import (
	commondomain "goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"

	"github.com/danielgtaylor/huma/v2"
)

// Creates a request id and attaches it to the request and context
func RequestIdMiddleware(ctx huma.Context, next func(huma.Context)) {
	requestId := commondomain.NewUUID()

	ctx.SetHeader("x-request-id", requestId.String())
	ctx = huma.WithValue(ctx, commoninfrastructure.RequestContextKey, commoninfrastructure.RequestContext{
		RequestId: requestId,
	})

	next(ctx)
}
