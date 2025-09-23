package sql

import (
	"context"
	"fmt"
	"goddd/internal/common"
	"goddd/internal/config"
	outboxdb "goddd/internal/infrastructure/sql/codegen/outbox"
	userdb "goddd/internal/infrastructure/sql/codegen/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"sql",
	fx.Provide(
		NewOutboxSql,
		NewUserSql,
		NewContext,
		NewDBPool,
	),
)

func NewOutboxSql(db common.DBTX) *outboxdb.Queries {
	return outboxdb.New(db)
}

func NewUserSql(db common.DBTX) *userdb.Queries {
	return userdb.New(db)
}

func NewContext() context.Context {
	return context.Background()
}

func NewDBPool(lc fx.Lifecycle, cfg *config.DBConfig, ctx context.Context) common.DBTX {
	pool, err := pgxpool.New(ctx,
		fmt.Sprintf(
			"host=%s dbname=%s user=%s password=%s",
			cfg.DBHost,
			cfg.DBName,
			cfg.DBUser,
			cfg.DBPassword,
		),
	)
	if err != nil {
		panic(err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Check connection
			return pool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})
	return pool
}
