package commonapplication

import (
	"context"
	"encoding/json"
	commondomain "goddd/internal/common/domain"
	"log/slog"
	"reflect"

	"go.uber.org/fx"
)

// Implementations
type EventHandler[T commondomain.DomainEventI] interface {
	Handle(ctx context.Context, log *slog.Logger, event T) error
}

// Generic EventHandler for fx
type EventHandlerInterface interface {
	HandlerEventType() commondomain.EventType
	Handle(ctx context.Context, payload []byte) error
}

// Adapter to fit EventHandler into EventHandlerInterface
type EventHandlerAdapter[T commondomain.DomainEventI] struct {
	log     *slog.Logger
	handler EventHandler[T]
}

func NewEventHandler[T commondomain.DomainEventI](log *slog.Logger, handler EventHandler[T]) *EventHandlerAdapter[T] {
	return &EventHandlerAdapter[T]{log: log, handler: handler}
}

func (e *EventHandlerAdapter[T]) HandlerEventType() commondomain.EventType {
	var t T
	return t.GetEventType()
}
func (e *EventHandlerAdapter[T]) Handle(ctx context.Context, payload []byte) error {
	// TODO; handlers do not respect ctx cancel
	// Wrap in a transaction ctx and check ctx.Done() before commit

	var event T
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return err
	}

	log := e.log.With("handler", reflect.TypeOf(e.handler).Elem().Name())
	log.InfoContext(ctx, "Event handler start")
	err = e.handler.Handle(ctx, log, event)
	if err != nil {
		log.ErrorContext(ctx, "Event handler error", "err", err)
	}

	return err
}

func AsEventHandler[A any, T commondomain.DomainEventI, H EventHandler[T]](constructor func(A) H) fx.Option {
	return fx.Options(
		fx.Provide(constructor),
		fx.Provide(fx.Annotate(
			func(log *slog.Logger, h H) EventHandlerInterface {
				return NewEventHandler(log, h)
			},
			fx.ResultTags(`group:"eventHandlers"`),
		)),
	)
}
