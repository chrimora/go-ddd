package eventhandlers

import (
	"goddd/internal/domain/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user_event_handlers",
	common.AsEventHandler(NewUserCreatedHandler, &UserCreatedHandler{}),
	common.AsEventHandler(NewUserCreatedHandler2, &UserCreatedHandler2{}),
)
