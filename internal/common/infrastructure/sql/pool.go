package commonsql

import (
	"context"
	"fmt"
	"goddd/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewContext() context.Context {
	return context.Background()
}

func NewDBPool(lc fx.Lifecycle, cfg *config.DBConfig, ctx context.Context) DBTX {
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
