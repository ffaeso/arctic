package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var partialOverrideConfigFile = filepath.Join("testdata", "partial_override.yml")

func TestLoad_ConfigFileOverride(t *testing.T) {
	cfg, err := Load(partialOverrideConfigFile)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestLoad_EnvOverride(t *testing.T) {
	// overrides default
	t.Setenv("ARCTIC_SERVER_PORT", "4200")
	// overrides config file value
	t.Setenv("ARCTIC_LOG_LEVEL", "error")

	cfg, err := Load(filepath.Join(partialOverrideConfigFile))
	require.NoError(t, err)

	assert.Equal(t, 4200, cfg.Server.Port)
	assert.Equal(t, "error", cfg.Log.Level)
}
