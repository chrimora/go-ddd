package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/outbox/application"
	"goddd/internal/user"

	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "worker"}
	fx.New(
		fx.Supply(service),
		common.Module,
		config.Module,
		outbox.WorkerModule,
		user.WorkerModule,
		fx.Invoke(func(*application.Worker) {}),
	).Run()
}
