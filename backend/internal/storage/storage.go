package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/pegov/enterboard/backend/internal/config"
)

type Storage struct {
	DB *sqlx.DB
}

func New(ctx context.Context, logger *slog.Logger, cfg *config.Config) (*Storage, error) {
	db, err := NewPG(ctx, logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("DB: %w", err)
	}

	return &Storage{DB: db}, nil
}
