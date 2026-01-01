package cliutils

import (
	"context"
	"fmt"

	"github.com/ffaeso/arctic/internal/config"
	"github.com/spf13/cobra"
)

// Prevent collisions with other packages that use
// the same keys in context as context keys are matched
// using type + value
type ContextKey string

const (
	ConfigKey ContextKey = "static_configs"
)

// GetConfig retrieves the Arctic static configuration from the command context.
// Returns an error if the config was not set via SetConfig or cannot be cast.
//
// This is mainly used in commands to get access to static configurations
// that were set at root command
func GetConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg, ok := cmd.Context().Value(ConfigKey).(*config.Config)
	if !ok {
		return nil, fmt.Errorf("unable to load config: failed to cast into *config.Config")
	}

	return cfg, nil
}

// SetConfig stores the Arctic configuration in the command context.
// This is typically called once in the root command's PersistentPreRunE.
func SetConfig(cmd *cobra.Command, cfg *config.Config) {
	ctx := context.WithValue(cmd.Context(), ConfigKey, cfg)
	cmd.SetContext(ctx)
}
