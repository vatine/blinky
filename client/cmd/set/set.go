package set

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/pkg/protos"
)

func setLEDs(cmd *cobra.Command, args []string) {
	config.Setup()

	log.Debug("pre-connect")
	client, err := config.Connect()
	log.Debug("post-connect")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get client")
		return
	}

	red, _ := strconv.Atoi(args[0])
	green, _ := strconv.Atoi(args[1])
	blue, _ := strconv.Atoi(args[2])
	req := protos.SetLEDRequest{
		Red:   int32(red),
		Green: int32(green),
		Blue:  int32(blue),
		LEDs:  config.LEDs,
	}

	_, err = client.SetLEDs(context.Background(), &req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"red":   red,
			"green": green,
			"blue":  blue,
		}).Error("SetLEDs failed")
	}
}

var Cmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(3),
	Use:   "set",
	Short: "set LEDs",
	Long:  "Set LEDs, takes red green blue as arguments",
	Run:   setLEDs,
}
