package application

import (
	"context"
	"goddd/internal/common/domain"
	"goddd/internal/config"
	"goddd/internal/outbox/domain"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"
)

type DomainEventForwarder struct {
	log        *slog.Logger
	txManager  *commondomain.TxManager
	outboxRepo domain.OutboxRepositoryI
	cfg        *config.ForwarderConfig
	publisher  EventPublisherI

	ticker *time.Ticker
	cancel context.CancelFunc
}

func NewForwarder(
	lc fx.Lifecycle,
	log *slog.Logger,
	txManager *commondomain.TxManager,
	outboxRepo domain.OutboxRepositoryI,
	cfg *config.ForwarderConfig,
	publisher EventPublisherI,
) *DomainEventForwarder {
	ctx, cancel := context.WithCancel(context.Background())

	forwarder := &DomainEventForwarder{
		log:        log,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		cfg:        cfg,
		publisher:  publisher,
		ticker:     time.NewTicker(cfg.WatchdogTick),
		cancel:     cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go forwarder.RunPublisher(ctx)
			go forwarder.RunWatchdog(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			forwarder.Stop()
			return nil
		},
	})
	return forwarder
}

func (w *DomainEventForwarder) RunPublisher(ctx context.Context) {
	w.log.Info("Forwarder publisher start")
	for {
		select {
		case <-ctx.Done():
			w.log.Info("Forwarder publisher stopped")
			return
		default:
			has_published := w.publishBatch(ctx)
			if !has_published {
				time.Sleep(w.cfg.PublisherSleep)
			}
		}
	}
}
func (w *DomainEventForwarder) RunWatchdog(ctx context.Context) {
	w.log.Info("Forwarder watchdog start")
	for {
		select {
		case <-ctx.Done():
			w.log.Info("Forwarder watchdog stopped")
			return
		case <-w.ticker.C:
			var count int
			t := time.Now().UTC().Add(-w.cfg.WatchdogStaleLimit)
			err := w.txManager.WithTx(ctx, func(tx pgx.Tx) error {
				c, err := w.outboxRepo.RequeueStaleEvents(ctx, tx, t, w.cfg.MaxRetries)
				count = c
				return err
			})
			if count > 0 {
				w.log.Warn("StaleEvent count", "count", count)
			}
			if err != nil {
				w.log.Error("StaleEvent requeue error", "err", err)
			}
		}
	}
}

func (w *DomainEventForwarder) Stop() {
	w.ticker.Stop()
	w.cancel()
}

func (w *DomainEventForwarder) publishBatch(ctx context.Context) bool {
	var events []*domain.OutboxEvent
	err := w.txManager.WithTx(ctx, func(tx pgx.Tx) error {
		e, err := w.outboxRepo.GetNextEventBatch(ctx, tx, w.cfg.PublisherBatchSize, w.cfg.MaxRetries)
		events = e
		return err
	})
	if err != nil {
		w.log.Error("EventBatch get error", "err", err)
		return false
	}
	if len(events) == 0 {
		return false
	}

	for _, event := range events {
		log := w.log.With("event", event)

		err := w.publisher.Publish(ctx, event)
		if err != nil {
			log.Error("Event publish error", "err", err)
		}

		err = w.txManager.WithTx(ctx, func(tx pgx.Tx) error {
			return w.outboxRepo.CompleteEvent(ctx, tx, event.ID)
		})
		if err != nil {
			log.Error("Event complete error", "err", err)
		}
	}
	return true
}
