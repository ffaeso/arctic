package main

import (
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx := t.Context()
	getenv := func(key string) string {
		return map[string]string{"PORT": "18080"}[key]
	}

	go run(ctx, getenv)
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/test")
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
	getenv := func(key string) string {
		return map[string]string{"PORT": "invalid"}[key]
	}

	err := run(ctx, getenv)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
	t.Logf("got expected error: %v", err)
}

func TestRunPortInUse(t *testing.T) {
	ctx := t.Context()
	getenv := func(key string) string {
		return map[string]string{"PORT": "18081"}[key]
	}

	go run(ctx, getenv)
	time.Sleep(100 * time.Millisecond)

	err := run(ctx, getenv)
	if err == nil {
		t.Fatal("expected error for port in use")
	}
	t.Logf("got expected error: %v", err)
}
