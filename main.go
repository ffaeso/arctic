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
)

// read the article below if confused on the structure
//
// source: https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/
func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "Arctic server stopped: %s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// load env, verify if any is missing, empty or not there

	// TODO: parse flags, load config, db pool, start server etc

	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := getenv("PORT")
	if port == "" {
		return fmt.Errorf("`PORT` environment variable doesn't exist")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
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

	log.Printf("Arctic server starting on port %s...", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
