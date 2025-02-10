package service

import (
	"context"
	"log/slog"

	"github.com/pegov/enterboard/backend/internal/dto"
	"github.com/pegov/enterboard/backend/internal/model"
	"github.com/pegov/enterboard/backend/internal/repo"
)

type Service struct {
	logger *slog.Logger
	r      *repo.Repo
}

func New(logger *slog.Logger, r *repo.Repo) *Service {
	return &Service{logger: logger, r: r}
}

func (s *Service) CreatePost(
	ctx context.Context,
	data dto.CreatePost,
) (*model.Post, error) {
	if err := data.Validate(); err != nil {
		return nil, dto.NewValidationError(err)
	}

	return s.r.CreatePost(ctx, data)
}
