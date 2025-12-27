package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/V4N1LLA-1CE/arctic/internal/config"
)

func TestRun(t *testing.T) {
	ctx := t.Context()

	cfg := &config.Config{
		Server: config.ServerConfig{Port: 18080},
		Postgres: config.PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			Name:     "test",
			SSLMode:  "disable",
		},
		Log: config.LogConfig{Level: "info", Format: "json"},
	}

	go run(ctx, cfg)
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/health")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRunInvalidPort(t *testing.T) {
	ctx := t.Context()

	cfg := &config.Config{
		Server:   config.ServerConfig{Port: -1},
		Postgres: config.PostgresConfig{},
		Log:      config.LogConfig{Level: "info", Format: "json"},
	}

	err := run(ctx, cfg)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
	t.Logf("got expected error: %v", err)
}

func TestRunPortInUse(t *testing.T) {
	ctx := t.Context()

	cfg := &config.Config{
		Server:   config.ServerConfig{Port: 18081},
		Postgres: config.PostgresConfig{},
		Log:      config.LogConfig{Level: "info", Format: "json"},
	}

	go run(ctx, cfg)
	time.Sleep(100 * time.Millisecond)

	err := run(ctx, cfg)
	if err == nil {
		t.Fatal("expected error for port in use")
	}
	t.Logf("got expected error: %v", err)
}

