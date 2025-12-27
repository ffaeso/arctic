package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPostgresConfig_DSN(t *testing.T) {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "arctic",
		Password: "secret",
		Name:     "arctic",
		SSLMode:  "disable",
	}

	expected := "postgres://arctic:secret@localhost:5432/arctic?sslmode=disable"
	if cfg.DSN() != expected {
		t.Errorf("expected %s, got %s", expected, cfg.DSN())
	}
}

func TestLoad(t *testing.T) {
	configPath := createTestConfig(t, `
server:
  port: 8080

postgres:
  host: localhost
  port: 5432
  user: arctic
  password: arctic
  name: arctic
  sslmode: disable

log:
  level: info
  format: json
`)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Postgres.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Postgres.Host)
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	configPath := createTestConfig(t, `
server:
  port: 8080

log:
  level: info
  format: json
`)

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected validation error")
	}
	t.Logf("got expected error: %v", err)
}

func createTestConfig(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	return configPath
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	t.Logf("got expected error: %v", err)
}

func TestLoad_InvalidYAML(t *testing.T) {
	configPath := createTestConfig(t, `
invalid: yaml: content: [[[
`)

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid yaml")
	}
	t.Logf("got expected error: %v", err)
}
