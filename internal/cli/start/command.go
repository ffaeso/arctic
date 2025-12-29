package start

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start Artic",
		Long:  "Start the Arctic server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting arctic...")
			// TODO: load and verify config -> testdb connection -> perform setup (db migrations) -> start server
			return nil
		},
	}
}
