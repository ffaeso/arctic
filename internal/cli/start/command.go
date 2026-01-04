package start

import (
	"github.com/ffaeso/arctic/internal/cliutils"
	"github.com/ffaeso/arctic/internal/database"
	"github.com/ffaeso/arctic/internal/logger"
	"github.com/ffaeso/arctic/internal/server"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start Arctic",
		Long:  "Start the Arctic server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := cliutils.GetConfig(cmd)
			if err != nil {
				return err
			}

			// instantiate logger
			l := logger.New(&cfg.Log)
			l.Info("loaded config", "config", cfg)

			// database connection
			conn, err := database.NewPostgresPool(&cfg.Datasource)
			if err != nil {
				return err
			}
			defer conn.Close()
			l.Info("successfully connected to database", "dsn", cfg.Datasource.DSN())

			// instantiate server
			srv := server.New(
				&cfg.Server,
				server.Mount(l),
				l,
			)

			// run server
			return srv.Serve()
		},
	}
}
