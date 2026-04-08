package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		common.ConsumerModule,
		config.ConsumerModule,
		outbox.ConsumerModule,
		// Order important - close forwarder first
		fx.Invoke(func(*outbox.Consumer) {}),
		fx.Invoke(func(*outbox.DomainEventForwarder) {}),
	).Run()
}
