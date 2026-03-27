package main

import (
	"goddd/internal/common"
	"goddd/internal/config"
	"goddd/internal/outbox"
	"goddd/internal/post"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		common.ConsumerModule,
		config.ConsumerModule,
		outbox.ConsumerModule,
		post.ConsumerModule,
		// Order important - close forwarder first
		fx.Invoke(func(*outbox.Consumer) {}),
		fx.Invoke(func(*outbox.DomainEventForwarder) {}),
	).Run()
}
