package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/pegov/enterboard/backend/internal/config"
)

type Storage struct {
	DB    *sqlx.DB
	Cache *redis.Client
}

func New(ctx context.Context, logger *slog.Logger, cfg *config.Config) (*Storage, error) {
	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)
	db, err := NewPG(
		ctx,
		logger,
		pgURL,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxOpenConns,
		cfg.DB.ConnMaxLifetime,
	)
	if err != nil {
		return nil, fmt.Errorf("pg: %w", err)
	}

	redisURL := fmt.Sprintf(
		"redis://%s:%d/%s",
		cfg.Cache.Host,
		cfg.Cache.Port,
		cfg.Cache.Database,
	)
	cache, err := NewRedis(ctx, logger, redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	return &Storage{DB: db, Cache: cache}, nil
}
