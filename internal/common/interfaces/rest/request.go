package commonrest

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type (
	// Request
	IdParam struct {
		ID uuid.UUID `path:"id"`
	}
	CreateRequest[T any] struct {
		Body T
	}
	UpdateRequest[T any] struct {
		IdParam
		Body T
	}

	// Response
	EmptyResponse   struct{}
	Response[T any] struct {
		Body T
	}

	// Payloads
	IdPayload struct {
		ID uuid.UUID `json:"id"`
	}
)

func BuildResponse[T any](body T) *Response[T] {
	return &Response[T]{Body: body}
}

func UnexpectedErrorResponse(log *slog.Logger, ctx context.Context, err error) error {
	log.ErrorContext(ctx, "REQUEST", "error", err.Error())
	return huma.Error500InternalServerError("Oops!")
}

func NotFoundResponse(log *slog.Logger, ctx context.Context, err error) error {
	log.WarnContext(ctx, "REQUEST", "error", err.Error())
	return huma.Error404NotFound(err.Error())
}
