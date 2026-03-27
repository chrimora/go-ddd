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
	return newDBPool(lc, cfg.WriteConnString(), cfg, ctx)
}

func NewReadDB(lc fx.Lifecycle, cfg *config.DBConfig, ctx context.Context) ReadDB {
	return newDBPool(lc, cfg.ReadConnString(), cfg, ctx)
}

func newDBPool(lc fx.Lifecycle, connString string, cfg *config.DBConfig, ctx context.Context) *pgxpool.Pool {
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}
	poolCfg.MaxConns = cfg.PoolMaxConns
	poolCfg.MinConns = cfg.PoolMinConns
	poolCfg.MaxConnIdleTime = cfg.PoolMaxConnIdle
	poolCfg.MaxConnLifetime = cfg.PoolMaxConnLife

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
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
