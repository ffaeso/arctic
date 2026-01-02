package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var partialOverrideConfigFile = filepath.Join("testdata", "partial_override.yml")

func TestLoad_ConfigFileOverride(t *testing.T) {
	clearEnvironmentVariables(t)

	cfg, err := Load(partialOverrideConfigFile)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestLoad_EnvOverride(t *testing.T) {
	clearEnvironmentVariables(t)

	// overrides default
	t.Setenv("ARCTIC_SERVER_ADDR", "4200")
	// overrides config file value
	t.Setenv("ARCTIC_LOG_LEVEL", "error")

	cfg, err := Load(filepath.Join(partialOverrideConfigFile))
	require.NoError(t, err)

	assert.Equal(t, 4200, cfg.Server.Addr)
	assert.Equal(t, "error", cfg.Log.Level)
}

func TestDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatasourceConfig
		expected string
	}{
		{
			name: "generic config",
			config: DatasourceConfig{
				Username: "arctic",
				Password: "changeme",
				Host:     "localhost",
				Port:     5554,
				DbName:   "arctic_dev_db",
				Sslmode:  "require",
			},
			expected: "postgres://arctic:changeme@localhost:5554/arctic_dev_db?sslmode=require",
		},
		{
			name: "different ssl mode",
			config: DatasourceConfig{
				Username: "user",
				Password: "pass",
				Host:     "db.example.com",
				Port:     5433,
				DbName:   "mydb",
				Sslmode:  "disable",
			},
			expected: "postgres://user:pass@db.example.com:5433/mydb?sslmode=disable",
		},
		{
			name: "special characters in password",
			config: DatasourceConfig{
				Username: "admin",
				Password: "p@ssw0rd!",
				Host:     "127.0.0.1",
				Port:     5432,
				DbName:   "prod",
				Sslmode:  "verify-full",
			},
			expected: "postgres://admin:p@ssw0rd!@127.0.0.1:5432/prod?sslmode=verify-full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.DSN()
			expected := tt.expected
			assert.Equal(t, expected, actual)
		})
	}
}

func clearEnvironmentVariables(t *testing.T) {
	// only clear environment variables that are optional / have defaults
	// so the tests don't break
	t.Helper()
	t.Setenv("ARCTIC_LOG_LEVEL", "")
	t.Setenv("ARCTIC_SERVER_ADDR", "")
}
