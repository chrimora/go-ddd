package commonsql

import (
	"context"
	"goddd/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewContext() context.Context {
	return context.Background()
}

func NewWriteDB(lc fx.Lifecycle, cfg *config.DBConfig, ctx context.Context) WriteDB {
	return newDBPool(lc, cfg.WriteConnString(), ctx)
}

func NewReadDB(lc fx.Lifecycle, cfg *config.DBConfig, ctx context.Context) ReadDB {
	return newDBPool(lc, cfg.WriteConnString(), ctx)
}

func newDBPool(lc fx.Lifecycle, connString string, ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connString)
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
