package start

import (
	"github.com/ffaeso/arctic/internal/config"
	"github.com/ffaeso/arctic/internal/logger"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start Arctic",
		Long:  "Start the Arctic server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// load config
			cfg, err := config.Load("")
			if err != nil {
				return err
			}

			// instantiate logger
			l := logger.New(cfg.Log)
			l.Info("loaded config", "config", cfg)

			return nil
		},
	}
}
