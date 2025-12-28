package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/V4N1LLA-1CE/arctic/internal/config"
	"github.com/spf13/viper"
)

// read the article below if confused on the structure
//
// source: https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/
func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %s\n", err)
		os.Exit(1)
	}

	log.Printf("config loaded from: %s", viper.ConfigFileUsed())
	log.Printf("config: %v", cfg)

	ctx := context.Background()
	if err := run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Arctic server stopped: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	// graceful shutdown - 30 seconds timeout
	go func() {
		// wait for context to finish
		<-ctx.Done()

		timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(timeout)
	}()

	log.Printf("Arctic server starting on port %d...", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
