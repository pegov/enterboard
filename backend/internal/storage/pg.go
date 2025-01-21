package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/pegov/enterboard/backend/internal/config"
)

func NewPG(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
) (*sqlx.DB, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)
	logger.Info("Parsing DB config...")
	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	logger.Info("Creating DB pool...")
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	logger.Info("Pinging DB...")
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	logger.Info("DB is online!")

	sqldb := stdlib.OpenDBFromPool(pool)
	sqldb.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqldb.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqldb.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	db := sqlx.NewDb(sqldb, "pgx")

	// TODO: migrations

	return db, nil
}
