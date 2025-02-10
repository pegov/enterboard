package dto

import (
	"errors"
	"strings"
)

type CreatePost struct {
	BoardID uint64 `json:"boardId"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

var (
	ErrTitleIsEmpty   = errors.New("title is empty")
	ErrMessageIsEmpty = errors.New("message is empty")
)

func (v *CreatePost) Validate() error {
	v.Title = strings.TrimSpace(v.Title)
	v.Message = strings.TrimSpace(v.Message)

	if v.Title == "" {
		return ErrTitleIsEmpty
	}

	if v.Message == "" {
		return ErrMessageIsEmpty
	}

	return nil
}
