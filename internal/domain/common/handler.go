package common

import (
	"context"
	"encoding/json"
	"log/slog"
	"reflect"

	"go.uber.org/fx"
)

// Implementations
type EventHandler[T DomainEventI] interface {
	Handle(ctx context.Context, event T) error
}

// Generic EventHandler for fx
type EventHandlerInterface interface {
	HandlerEventType() EventType
	Handle(ctx context.Context, payload []byte) error
}

// Adapter to fit EventHandler into EventHandlerInterface
type EventHandlerAdapter[T DomainEventI] struct {
	log     *slog.Logger
	handler EventHandler[T]
}

func NewEventHandler[T DomainEventI](log *slog.Logger, handler EventHandler[T]) *EventHandlerAdapter[T] {
	return &EventHandlerAdapter[T]{log: log, handler: handler}
}

func (e *EventHandlerAdapter[T]) HandlerEventType() EventType {
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

	handlerName := reflect.TypeOf(e.handler).Elem().Name()
	e.log.InfoContext(ctx, "Event handler start", "handler", handlerName)
	err = e.handler.Handle(ctx, event)
	if err != nil {
		e.log.ErrorContext(ctx, "Event handler error", "handler", handlerName, "err", err)
	}

	return err
}

// _type is a hack to infer types
func AsEventHandler[T DomainEventI, H EventHandler[T]](constructor any, _type H) fx.Option {
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
