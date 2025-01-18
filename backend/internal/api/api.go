package api

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pegov/enterboard/backend/internal/http/bind"
	"github.com/pegov/enterboard/backend/internal/http/render"
)

func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	handler := NewHandler(logger)
	makeHandler := func(fn HandlerFuncWithError) http.HandlerFunc {
		return makeHandlerFull(fn, logger)
	}

	apiV1 := chi.NewRouter()
	apiV1.Post("/posts", makeHandler(handler.CreatePost))

	r.Mount("/api/v1", apiV1)

	logger.Debug("Listen")
	http.ListenAndServe(":3000", r)
}

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

type HandlerFuncWithError = func(w http.ResponseWriter, r *http.Request) error

type Detail struct {
	Detail string `json:"detail"`
}

func NewDetail(detail string) Detail {
	return Detail{
		Detail: detail,
	}
}

func makeHandlerFull(fn HandlerFuncWithError, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bindJSONError *bind.BindJSONError
		if err := fn(w, r); err != nil {
			switch {
			case errors.As(err, &bindJSONError):
				render.JSON(
					w,
					bindJSONError.Status,
					NewDetail(bindJSONError.Message),
				)

			default:
				logger.Error("Internal server error", slog.Any("err", err))
				render.String(
					w,
					http.StatusInternalServerError,
					"Internal server error",
				)
			}
		}
	}
}

type CreatePost struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) error {
	var req CreatePost
	if err := bind.JSON(r, &req); err != nil {
		return err
	}

	return nil
}
