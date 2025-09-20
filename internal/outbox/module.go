package outbox

import (
	"gotemplate/internal/common"
	outboxdb "gotemplate/internal/infrastructure/sql/codegen/outbox"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbox",
	fx.Provide(
		NewOutboxRepository,
		NewSql,
	),
)

func NewSql(db common.DBTX) *outboxdb.Queries {
	return outboxdb.New(db)
}
