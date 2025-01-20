package service

import (
	"context"
	"log/slog"

	"github.com/pegov/enterboard/backend/internal/model"
)

type Service struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) CreatePost(
	ctx context.Context,
	data model.CreatePost,
) (*model.Post, error) {
	return nil, nil
}
