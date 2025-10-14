package application

import (
	"context"
	"goddd/internal/common/application"
	"goddd/internal/common/infrastructure"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/fx"
)

type Consumer struct{}

func NewConsumer(
	handlers []commonapplication.EventHandlerInterface,
	lc fx.Lifecycle,
	log *slog.Logger,
	router *message.Router,
	subscriber message.Subscriber,
) *Consumer {
	// Register handlers
	for _, handler := range handlers {
		router.AddConsumerHandler(
			handler.GetName(),
			string(handler.GetType()),
			subscriber,
			func(msg *message.Message) error {
				trace := infrastructure.NewTraceCtxFromMessage(msg.Metadata)
				return handler.Handle(trace.ToCtx(context.Background()), msg.UUID, msg.Payload)
			},
		)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("Consumer start")
			go func() {
				err := router.Run(context.Background())
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Consumer stopping")
			err := router.Close()
			return err
		},
	})
	return &Consumer{}
}
