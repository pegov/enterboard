package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPG(
	ctx context.Context,
	logger *slog.Logger,
	url string,
	maxIdleConns int,
	maxOpenConns int,
	connMaxLifetime time.Duration,
) (*sqlx.DB, error) {
	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}
	poolCfg.ConnConfig.Tracer = &customQueryTracer{
		log: logger,
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	sqldb := stdlib.OpenDBFromPool(pool)
	sqldb.SetMaxIdleConns(maxIdleConns)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetConnMaxLifetime(connMaxLifetime)

	db := sqlx.NewDb(sqldb, "pgx")

	initScript, err := os.ReadFile("resources/sql/init.sql")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	if _, err := sqldb.ExecContext(ctx, string(initScript)); err != nil {
		return nil, fmt.Errorf("sql.Exec initScript: %w", err)
	}

	return db, nil
}

type customQueryTracer struct {
	log *slog.Logger
}

func (tracer *customQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	ctx = context.WithValue(ctx, "query", data.SQL)
	ctx = context.WithValue(ctx, "args", data.Args)
	ctx = context.WithValue(ctx, "start", time.Now())
	return ctx
}

func (tracer *customQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		tracer.log.Error(
			"pgx",
			slog.Any("query", ctx.Value("query")),
			slog.Any("args", ctx.Value("args")),
			slog.Any("duration", time.Since(ctx.Value("start").(time.Time))),
			slog.Any("err", data.Err),
		)
	}
}
