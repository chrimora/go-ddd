package outbox

import (
	"context"
	"fmt"
	"goddd/internal/common"
	domain "goddd/internal/domain/common"
	outboxdb "goddd/internal/infrastructure/sql/codegen/outbox"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type Dispatcher struct {
	log      *slog.Logger
	handlers map[domain.EventType][]domain.EventHandlerInterface
}

func NewDispatcher(handlers []domain.EventHandlerInterface, log *slog.Logger) *Dispatcher {
	handlersMap := make(map[domain.EventType][]domain.EventHandlerInterface)

	for _, handler := range handlers {
		eventType := handler.HandlerEventType()
		handlersMap[eventType] = append(handlersMap[eventType], handler)
	}

	return &Dispatcher{
		log:      log,
		handlers: handlersMap,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event *outboxdb.EventOutbox) error {
	handlers, ok := d.handlers[domain.EventType(event.EventType)]
	if !ok {
		return fmt.Errorf("no handler for %s", event.EventType)
	}

	// errgroup will fail fast - cancel all handlers on first failure
	group, groupCtx := errgroup.WithContext(ctx)
	groupCtx = common.JsonToCtx(event.EventContext, groupCtx)

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
