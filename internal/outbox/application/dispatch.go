package application

import (
	"context"
	"fmt"
	"goddd/internal/common/application"
	"goddd/internal/common/domain"
	"goddd/internal/common/infrastructure"
	"goddd/internal/outbox/infrastructure/sql"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type Dispatcher struct {
	log      *slog.Logger
	handlers map[commondomain.EventType][]commonapplication.EventHandlerInterface
}

func NewDispatcher(
	handlers []commonapplication.EventHandlerInterface,
	log *slog.Logger,
) *Dispatcher {
	handlersMap := make(map[commondomain.EventType][]commonapplication.EventHandlerInterface)

	for _, handler := range handlers {
		eventType := handler.HandlerEventType()
		handlersMap[eventType] = append(handlersMap[eventType], handler)
	}

	return &Dispatcher{
		log:      log,
		handlers: handlersMap,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event *sql.EventOutbox) error {
	handlers, ok := d.handlers[commondomain.EventType(event.EventType)]
	if !ok {
		return fmt.Errorf("no handler for %s", event.EventType)
	}

	// errgroup will fail fast - cancel all handlers on first failure
	group, groupCtx := errgroup.WithContext(ctx)
	groupCtx = infrastructure.JsonToCtx(event.EventContext, groupCtx)

	start := time.Now()
	log := d.log.With("event", event.ID, "type", event.EventType)

	log.InfoContext(groupCtx, "Event handling start")
	for _, handler := range handlers {
		h := handler // avoid loop capture

		group.Go(func() error {
			return h.Handle(groupCtx, event.Payload)
		})
	}

	err := group.Wait()
	log.InfoContext(groupCtx, "Event handling executed", "duration", time.Since(start))
	return err
}
