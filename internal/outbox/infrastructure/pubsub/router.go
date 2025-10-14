package pubsub

import (
	"goddd/internal/config"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func NewRouter(log *slog.Logger, cfg *config.RouterConfig) *message.Router {
	logger := watermill.NewSlogLogger(log)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.Retry{
			MaxRetries:      cfg.MaxRetries,
			InitialInterval: cfg.RetryInterval,
			Multiplier:      cfg.RetryIntervalMultiplier,
			Logger:          logger,
		}.Middleware,
	)

	return router
}
