package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/user"

	"go.uber.org/fx"
)

func main() {
	service := config.ServiceConfig{Name: "consumer"}
	fx.New(
		fx.Supply(service),
		common.ConsumerModule,
		config.ConsumerModule,
		outbox.ConsumerModule,
		user.ConsumerModule,
		// Order important - close forwarder first
		fx.Invoke(func(*outbox.Consumer) {}),
		fx.Invoke(func(*outbox.DomainEventForwarder) {}),
	).Run()
}
