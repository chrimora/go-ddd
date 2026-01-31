package middleware

import (
	"goddd/internal/common/infrastructure"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Creates a request id and attaches it to the request and context
func RequestIdMiddleware(ctx huma.Context, next func(huma.Context)) {
	requestId := uuid.NewString()

	ctx.SetHeader("x-request-id", requestId)
	ctx = huma.WithValue(ctx, commoninfrastructure.RequestIdKey, requestId)

	next(ctx)
}
