package logger

import (
	"log/slog"
	"os"

	"github.com/fraeso/arctic/internal/config"
)

// New creates and returns a new *slog.Logger instance
// with settings based on the loggerConfig passed in.
//
// At the moment this supports log formats [text | json]
// and logging level [debug | info | warn | error]
func New(loggerConfig *config.LogConfig) *slog.Logger {
	// convert log level text from config to slog.Level
	var level slog.Level
	level.UnmarshalText([]byte(loggerConfig.Level))

	opts := slog.HandlerOptions{Level: level}

	if loggerConfig.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &opts))
}
