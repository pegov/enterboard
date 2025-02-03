package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/pegov/enterboard/backend/internal/api"
	"github.com/pegov/enterboard/backend/internal/config"
	"github.com/pegov/enterboard/backend/internal/repo"
	"github.com/pegov/enterboard/backend/internal/service"
	"github.com/pegov/enterboard/backend/internal/storage"
	"github.com/pegov/enterboard/backend/internal/util"
)

var (
	verbose = flag.Bool("verbose", false, "log level = debug")
)

func main() {
	flag.Parse()
	godotenv.Load()

	cfg := config.New()

	var level slog.Level
	if *verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	logger := slog.New(util.NewColoredHandler(os.Stdout, &util.ColoredHandlerOptions{
		Level: level,
	}))

	ctx := context.TODO()
	st, err := storage.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect to storages", slog.Any("err", err))
		os.Exit(1)
	}
	r := repo.New(st)

	srv := service.New(logger, r)

	api.Run(ctx, logger, cfg, srv)
}
