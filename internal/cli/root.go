package cli

import (
	"os"

	"github.com/fraeso/arctic/internal/cli/start"
	"github.com/fraeso/arctic/internal/cliutils"
	"github.com/fraeso/arctic/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arctic",
	Short: "Open-Source Identity and Access Management",
	Long:  `Arctic is an Open-Source Identity and Access Management (IAM) platform`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// load config
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		cliutils.SetConfig(cmd, cfg)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(start.NewCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
