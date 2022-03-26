package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/cmd/get"
	"github.com/vatine/blinky/client/cmd/set"
)

var rootCmd = &cobra.Command{
	Short: "Client for the blinky system",
}

func Execute() error {
	rootCmd.PersistentFlags().StringVar(&config.Server, "server", "192.168.1.227:4004", "Address and port of the server.")
	rootCmd.PersistentFlags().StringVar(&config.LogLevel, "loglevel", "info", "Log level")
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(set.Cmd)

	return rootCmd.Execute()
}
