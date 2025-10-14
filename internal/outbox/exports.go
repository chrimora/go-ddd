package outbox

import (
	"goddd/internal/outbox/application"
	"goddd/internal/outbox/domain"
)

type (
	DomainEventForwarder = application.DomainEventForwarder
	Consumer             = application.Consumer
	OutboxRepositoryI    = domain.OutboxRepositoryI
)
