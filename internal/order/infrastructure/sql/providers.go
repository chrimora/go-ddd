package sql

import commonsql "goddd/internal/common/infrastructure/sql"

type WriteOrderSql struct{ *Queries }
type ReadOrderSql struct{ *Queries }

func NewWriteOrderSql(db commonsql.WriteDB) WriteOrderSql {
	return WriteOrderSql{New(db)}
}
func NewReadOrderSql(db commonsql.ReadDB) ReadOrderSql {
	return ReadOrderSql{New(db)}
}
