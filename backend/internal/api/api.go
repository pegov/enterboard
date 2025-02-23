package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/pegov/enterboard/backend/internal/config"
	"github.com/pegov/enterboard/backend/internal/dto"
	"github.com/pegov/enterboard/backend/internal/http/bind"
	"github.com/pegov/enterboard/backend/internal/http/render"
	"github.com/pegov/enterboard/backend/internal/service"
)

func Run(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
	srv *service.Service,
) {
	reg := prometheus.NewRegistry()

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler := NewHandler(logger, srv)
	makeHandler := func(fn HandlerFuncWithError) http.HandlerFunc {
		return makeHandlerFull(fn, logger)
	}

	mMux := http.NewServeMux()
	mMux.Handle("/metrics", promHandler)

	apiV1 := chi.NewRouter()
	apiV1.Post("/posts", makeHandler(handler.CreatePost))

	r.Mount("/api/v1", apiV1)
	logger = logger.With(slog.String("unit", "[API]"))

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.App.MetricsHost, cfg.App.MetricsPort)
		logger.Debug("Metrics", slog.Any("addr", addr))
		http.ListenAndServe(addr, mMux)
	}()

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	logger.Debug("API", slog.Any("addr", addr))
	http.ListenAndServe(addr, r)
}

type Handler struct {
	logger *slog.Logger
	s      *service.Service
}

func NewHandler(logger *slog.Logger, s *service.Service) *Handler {
	return &Handler{logger: logger, s: s}
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
		var validationError *dto.ValidationError
		if err := fn(w, r); err != nil {
			switch {
			case errors.As(err, &bindJSONError):
				render.JSON(
					w,
					bindJSONError.Status,
					NewDetail(bindJSONError.Message),
				)

			case errors.As(err, &validationError):
				render.JSON(
					w,
					http.StatusBadRequest,
					NewDetail(validationError.Error()),
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

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) error {
	return GenericReqAndRes(w, r, h.s.CreatePost, "CreatePost")
}

func GenericReqAndRes[REQ, RES any](
	w http.ResponseWriter,
	r *http.Request,
	action func(context.Context, REQ) (*RES, error),
	actionName string,
) error {
	var req REQ
	if err := bind.JSON(r, &req); err != nil {
		return err
	}

	v, err := action(r.Context(), req)
	if err != nil {
		return fmt.Errorf("service.%s: %w", actionName, err)
	}

	render.JSON(w, http.StatusOK, v)
	return nil
}
