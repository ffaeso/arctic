package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ffaeso/arctic/internal/config"
)

type Server struct {
	// server port
	Addr int

	// multiplexer and loggers
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

	// start listener in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("unexpected server error", "error", err)
		}
	}()

	// perform graceful shutdown with 30 seconds timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Logger.Info("performing graceful shutdown...")
	return srv.Shutdown(ctx)
}
