package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/user"

	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "worker"}
	fx.New(
		fx.Supply(service),
		common.CoreModule,
		config.Module,
		outbox.WorkerModule,
		user.WorkerModule,
		fx.Invoke(func(*outbox.Worker) {}),
	).Run()
}
