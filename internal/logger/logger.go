package logger

import (
	"log/slog"
	"os"

	"github.com/ffaeso/arctic/internal/config"
)

func New(loggerConfig *config.LogConfig) *slog.Logger {
	// change log level text from config to slog.Level
	var level slog.Level
	level.UnmarshalText([]byte(loggerConfig.Level))

	opts := slog.HandlerOptions{Level: level}

	if loggerConfig.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &opts))
}
