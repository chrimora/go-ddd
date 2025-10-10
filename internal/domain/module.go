package domain

import (
	"goddd/internal/domain/common"
	"goddd/internal/domain/outbox"
	"goddd/internal/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"domain",
	fx.Provide(
		common.NewTxFactory,
		common.NewTxManager,
	),
	outbox.Module,
	user.Module,
)

type (
	EventType    = common.EventType
	DomainEventI = common.DomainEventI

	User            = user.User
	UserServiceI    = user.UserServiceI
	UserRepositoryI = user.UserRepositoryI
)
