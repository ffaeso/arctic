package cli

import (
	"os"

	"github.com/ffaeso/arctic/internal/cli/start"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arctic",
	Short: "Open-Source Identity and Access Management",
	Long:  `Arctic is an Open-Source Identity and Access Management (IAM) platform`,
}

func init() {
	rootCmd.AddCommand(start.NewCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
