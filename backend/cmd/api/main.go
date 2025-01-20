package main

import (
	"github.com/joho/godotenv"

	"github.com/pegov/enterboard/backend/internal/api"
	"github.com/pegov/enterboard/backend/internal/config"
)

func main() {
	godotenv.Load()

	cfg := config.New()
	api.Run(cfg)
}
