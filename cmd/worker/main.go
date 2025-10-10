package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/domain"
	eventhandlers "goddd/internal/domain/user/event_handlers"
	"goddd/internal/infrastructure/outbox"
	"goddd/internal/infrastructure/sql"

	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "worker"}
	fx.New(
		fx.Supply(service),
		fx.Provide(
			common.NewLogger,
			fx.Annotate(outbox.NewDispatcher, fx.ParamTags(`group:"eventHandlers"`)),
			outbox.NewWorker,
		),
		config.Module,
		domain.Module,
		sql.Module,
		eventhandlers.Module,
		fx.Invoke(func(*outbox.Worker) {}),
	).Run()
}
