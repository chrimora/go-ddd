package outbox

import (
	"goddd/internal/outbox/application"
	"goddd/internal/outbox/domain"
)

type (
	Worker            = application.Worker
	OutboxRepositoryI = domain.OutboxRepositoryI
)
