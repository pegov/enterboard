package repo

import (
	"context"

	"github.com/pegov/enterboard/backend/internal/model"
	"github.com/pegov/enterboard/backend/internal/storage"
)

type Repo struct {
	s *storage.Storage
}

func New(s *storage.Storage) *Repo {
	return &Repo{s: s}
}

func (r *Repo) CreatePost(ctx context.Context, data model.CreatePost) (*model.Post, error) {
	return nil, nil
}
