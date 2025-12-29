package start

import (
	"log"

	"github.com/ffaeso/arctic/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start Arctic",
		Long:  "Start the Arctic server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("")
			if err != nil {
				return err
			}

			log.Printf("%+v", cfg)

			return nil
		},
	}
}
