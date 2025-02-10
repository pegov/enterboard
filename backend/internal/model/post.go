package model

type Post struct {
	ID      uint64 `json:"id"`
	BoardID uint64 `json:"boardId"`
	Message string `json:"message"`
}
