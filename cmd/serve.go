package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/theopenlane/riverboat/internal/river"
	"github.com/theopenlane/riverboat/internal/server/serveropts"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the riverboat server",
	Run: func(cmd *cobra.Command, args []string) {
		err := serve(cmd.Context())
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().String("config", "./config/.config.yaml", "config file location")
}

func serve(ctx context.Context) error {
	serverOpts := []serveropts.ServerOption{}

	so := serveropts.NewServerOptions(serverOpts, k.String("config"))

	// pass the logger options to the job queue
	so.Config.Settings.JobQueue.Logger.Debug = k.Bool(debugFlag)
	so.Config.Settings.JobQueue.Logger.Pretty = k.Bool(prettyFlag)

	return river.Start(ctx, so.Config.Settings.JobQueue)
}
