package server

import (
	"log/slog"
	"net/http"

	"github.com/ffaeso/arctic/internal/api/http/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

func Mount(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(logger, nil))
	r.Use(middleware.Recoverer)

	r.Get("/ping", health.Ping)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
	})

	return r
}
