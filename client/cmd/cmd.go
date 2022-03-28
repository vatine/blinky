package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/cmd/get"
	"github.com/vatine/blinky/client/cmd/runner"
	"github.com/vatine/blinky/client/cmd/set"
)

var rootCmd = &cobra.Command{
	Short: "Client for the blinky system",
}

func Execute() error {
	rootCmd.PersistentFlags().StringVar(&config.Server, "server", "192.168.1.227:4004", "Address and port of the server.")
	rootCmd.PersistentFlags().StringVar(&config.LogLevel, "loglevel", "info", "Log level")
	rootCmd.PersistentFlags().StringVar(&config.LEDString, "leds", "", "LEDs to act on.")
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(set.Cmd)
	rootCmd.AddCommand(&runner.Cmd)

	return rootCmd.Execute()
}
