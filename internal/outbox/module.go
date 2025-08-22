package outbox

import (
	"gotemplate/internal/common"
	outboxdb "gotemplate/internal/outbox/db"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		NewOutboxRepository,
		NewSql,
	),
)

func NewSql(db common.DBTX) *outboxdb.Queries {
	return outboxdb.New(db)
}
