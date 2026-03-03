package sql

import commonsql "goddd/internal/common/infrastructure/sql"

type WriteUserSql struct{ *Queries }
type ReadUserSql struct{ *Queries }

func NewWriteUserSql(db commonsql.WriteDB) WriteUserSql {
	return WriteUserSql{New(db)}
}
func NewReadUserSql(db commonsql.ReadDB) ReadUserSql {
	return ReadUserSql{New(db)}
}
