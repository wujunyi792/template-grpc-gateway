package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"pinnacle-primary-be/cmd/config"
	"pinnacle-primary-be/cmd/create"
	"pinnacle-primary-be/cmd/server"
)

var rootCmd = &cobra.Command{
	Use:          "app",
	Short:        "app",
	SilenceUsage: true,
	Long:         `app`,
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
