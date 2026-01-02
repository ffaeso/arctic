//go:build integration

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ffaeso/arctic/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var partialOverrideConfigFile = filepath.Join("testdata", "partial_override.yml")

func TestMain(m *testing.M) {
	// clear all ARCTIC_* environment variables for complete isolation
	testutil.ClearArcticEnvVars()

	// set required datasource defaults for ALL tests
	testutil.SetupTestDatasource()

	code := m.Run()
	os.Exit(code)
}

func TestLoad_ConfigFileOverride(t *testing.T) {
	cfg, err := Load(partialOverrideConfigFile)
	require.NoError(t, err)

	assert.Equal(t, "debug", cfg.Log.Level)
}

func TestLoad_EnvOverride(t *testing.T) {
	// overrides default
	t.Setenv("ARCTIC_SERVER_ADDR", "4200")
	// overrides config file value
	t.Setenv("ARCTIC_LOG_LEVEL", "error")

	cfg, err := Load(filepath.Join(partialOverrideConfigFile))
	require.NoError(t, err)

	assert.Equal(t, 4200, cfg.Server.Addr)
	assert.Equal(t, "error", cfg.Log.Level)
}

func TestLoad_Precedence_AllThreeSources(t *testing.T) {
	// config file has: log.level = "debug"
	// default has: log.level = "info"
	// env var sets: log.level = "warn"
	t.Setenv("ARCTIC_LOG_LEVEL", "warn")

	cfg, err := Load(partialOverrideConfigFile)
	require.NoError(t, err)

	// env should win over config file and defaults
	assert.Equal(t, "warn", cfg.Log.Level)
}

func TestLoad_Precedence_FileAndDefaults(t *testing.T) {
	cfg, err := Load(partialOverrideConfigFile)
	require.NoError(t, err)

	// config file has: log.level = "debug"
	// default has: server.addr = 9726
	assert.Equal(t, "debug", cfg.Log.Level)           // file wins over default
	assert.Equal(t, DefaultHttpAddr, cfg.Server.Addr) // default used when not in file
}

func TestLoad_Precedence_OnlyDefaults(t *testing.T) {
	// load with empty/non-existent config file
	cfg, err := Load("")
	require.NoError(t, err)

	// all values should be defaults
	assert.Equal(t, DefaultLogLevel, cfg.Log.Level)
	assert.Equal(t, DefaultLogFormat, cfg.Log.Format)
	assert.Equal(t, DefaultHttpAddr, cfg.Server.Addr)
}

func TestLoad_ConfigFileNotFound_ExplicitPath(t *testing.T) {
	// when explicit non-existent path is provided, should return error
	cfg, err := Load("testdata/does_not_exist.yml")
	require.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_ConfigFileNotFound_SearchPaths(t *testing.T) {
	// when no path provided and config not found in search paths,
	// should not error and use defaults + env vars
	cfg, err := Load("")
	require.NoError(t, err)

	// should fall back to defaults
	assert.Equal(t, DefaultLogLevel, cfg.Log.Level)
	assert.Equal(t, DefaultHttpAddr, cfg.Server.Addr)
}

func TestLoad_MalformedYAML(t *testing.T) {
	// malformed YAML should return an error
	cfg, err := Load("testdata/malformed.yml")
	require.Error(t, err)
	assert.Nil(t, cfg)
}
