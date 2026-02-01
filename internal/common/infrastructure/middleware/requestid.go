package middleware

import (
	commondomain "goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"

	"github.com/danielgtaylor/huma/v2"
)

// Creates a request id and attaches it to the request and context
func RequestIdMiddleware(ctx huma.Context, next func(huma.Context)) {
	requestId := commondomain.NewUUID().String()

	ctx.SetHeader("x-request-id", requestId)
	ctx = huma.WithValue(ctx, commoninfrastructure.RequestIdKey, requestId)

	next(ctx)
}
