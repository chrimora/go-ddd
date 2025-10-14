package pubsub

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func NewGoChannel(log *slog.Logger) *gochannel.GoChannel {
	return gochannel.NewGoChannel(
		gochannel.Config{
			// Reduces in memory message to one
			// Reduces loss given unexpected close
			// Will have scaling issues
			BlockPublishUntilSubscriberAck: true,
		},
		watermill.NewSlogLogger(log),
	)
}
