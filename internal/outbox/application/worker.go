package application

import (
	"context"
	"errors"
	common "goddd/internal/common/domain"
	"goddd/internal/outbox/domain"
	outboxdb "goddd/internal/outbox/infrastructure/sql/codegen"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

type Worker struct {
	log        *slog.Logger
	txManager  *common.TxManager
	outboxRepo domain.OutboxRepositoryI
	dispatch   *Dispatcher

	ticker *time.Ticker
	cancel context.CancelFunc
}

func NewWorker(
	lc fx.Lifecycle,
	log *slog.Logger,
	txManager *common.TxManager,
	outboxRepo domain.OutboxRepositoryI,
	dispatch *Dispatcher,
) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	worker := &Worker{
		log:        log,
		txManager:  txManager,
		outboxRepo: outboxRepo,
		dispatch:   dispatch,
		ticker:     time.NewTicker(500 * time.Millisecond),
		cancel:     cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go worker.Run(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			worker.Stop()
			return nil
		},
	})
	return worker
}

func (w *Worker) Run(ctx context.Context) {
	w.log.Info("Worker start")
	for {
		select {
		case <-w.ticker.C:
			// New context to avoid cancelling handlers from service context
			w.processNext(context.Background())
		case <-ctx.Done():
			// Cancelled
			return
		}
	}
}

func (w *Worker) Stop() {
	w.log.Info("Worker stopping")
	w.ticker.Stop()
	w.cancel()
}

func (w *Worker) processNext(ctx context.Context) {
	var event *outboxdb.EventOutbox
	err := w.txManager.WithTxCtx(ctx, func(txCtx context.Context) error {
		e, err := w.outboxRepo.GetNextEvent(txCtx)
		event = e
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			// Backoff?
			return
		default:
			w.log.Error("Event get error", "err", err)
			return
		}
	}
	log := w.log.With("event", event.ID)

	err = w.dispatch.Dispatch(ctx, event)
	if err != nil {
		log.Info("Event requeue")
		err := w.outboxRepo.RequeueEvent(ctx, event.ID)
		if err != nil {
			log.Error("Event requeue error", "err", err)
			// TODO; stale event
		}
		return
	}

	err = w.outboxRepo.CompleteEvent(ctx, event.ID)
	if err != nil {
		log.Error("Error completing event", "err", err)
	}
}
