package application

import (
	"context"
	"goddd/internal/common/domain"
	"goddd/internal/outbox/domain"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

const RetriesCutoff = 5
const PublisherBatchSize = 10
const PublisherBackoff = 300 * time.Millisecond
const WatchdogTick = 3 * time.Second
const WatchdogStaleLimit = 3 * time.Second

type Worker struct {
	log        *slog.Logger
	txManager  *commondomain.TxManager
	outboxRepo domain.OutboxRepositoryI
	dispatch   *Dispatcher

	ticker *time.Ticker
	cancel context.CancelFunc
}

func NewWorker(
	lc fx.Lifecycle,
	log *slog.Logger,
	txManager *commondomain.TxManager,
	outboxRepo domain.OutboxRepositoryI,
	dispatch *Dispatcher,
) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	worker := &Worker{
		log:        log,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		dispatch:   dispatch,
		ticker:     time.NewTicker(WatchdogTick),
		cancel:     cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go worker.RunPublisher(ctx)
			go worker.RunWatchdog(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			worker.Stop()
			return nil
		},
	})
	return worker
}

func (w *Worker) RunPublisher(ctx context.Context) {
	w.log.Info("Worker publisher start")
	for {
		select {
		case <-ctx.Done():
			// Cancelled
			return
		default:
			has_published := w.publishBatch(ctx)
			if !has_published {
				time.Sleep(PublisherBackoff)
			}
		}
	}
}
func (w *Worker) RunWatchdog(ctx context.Context) {
	w.log.Info("Worker watchdog start")
	for {
		select {
		case <-ctx.Done():
			// Cancelled
			return
		case <-w.ticker.C:
			t := time.Now().UTC().Add(-WatchdogStaleLimit)
			count, err := w.outboxRepo.RequeueStaleEvents(ctx, t, RetriesCutoff)
			if count > 0 {
				w.log.Warn("StaleEvents", "count", count)
			}
			if err != nil {
				w.log.Error("StaleEvent requeue error", "err", err)
			}
		}
	}
}

func (w *Worker) Stop() {
	w.log.Info("Worker stopping")
	w.ticker.Stop()
	w.cancel()
}

func (w *Worker) publishBatch(ctx context.Context) bool {
	var events []*domain.OutboxEvent
	err := w.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		e, err := w.outboxRepo.GetNextEventBatch(txCtx, PublisherBatchSize, RetriesCutoff)
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

		// TODO; publish

		err = w.outboxRepo.CompleteEvent(ctx, event.ID)
		if err != nil {
			log.Error("Error completing event", "err", err)
		}
	}
	return true
}
