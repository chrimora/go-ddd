package sql

import commonsql "goddd/internal/common/infrastructure/sql"

type WritePostSql struct{ *Queries }
type ReadPostSql struct{ *Queries }

func NewWritePostSql(db commonsql.WriteDB) WritePostSql {
	return WritePostSql{New(db)}
}
func NewReadPostSql(db commonsql.ReadDB) ReadPostSql {
	return ReadPostSql{New(db)}
}
