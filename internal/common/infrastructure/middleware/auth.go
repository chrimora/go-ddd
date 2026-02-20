package middleware

import (
	"goddd/internal/common/infrastructure"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// Authenticate the user and attach their id to the context
func UserAuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	// TODO; replace this implementation

	userId := uuid.Nil.String()

	ctx = huma.WithValue(ctx, commoninfrastructure.UserIdKey, userId)

	next(ctx)
}
