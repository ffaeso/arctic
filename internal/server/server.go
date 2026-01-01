package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ffaeso/arctic/internal/config"
)

type Server struct {
	Addr int

	Mux    http.Handler
	Logger *slog.Logger
}

func New(
	cfg *config.ServerConfig,
	mux http.Handler,
	logger *slog.Logger,
) *Server {
	return &Server{
		Addr:   cfg.Addr,
		Mux:    mux,
		Logger: logger,
	}
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Addr),
		Handler: s.Mux,

		// TODO: configure more server settings
	}

	s.Logger.Info("starting arctic server...", "addr", s.Addr)
	return srv.ListenAndServe()
}
