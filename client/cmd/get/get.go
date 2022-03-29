package get

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/pkg/protos"
)

func getLEDs(cms *cobra.Command, args []string) {
	config.Setup()
	log.Debug("getLEDs called")
	req := protos.GetLEDRequest{LEDs: config.LEDs}

	log.Debug("pre-connect")
	client, err := config.Connect()
	log.Debug("post-connect")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get client")
		return
	}

	log.Debug("pre-call")
	resp, err := client.GetLEDs(context.Background(), &req)
	log.Debug("post-call")

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error gettling LED status")
		return
	}

	for _, status := range resp.Status {
		fmt.Printf("LED: %d, %d/%d/%d\n", status.GetLED(), status.GetRed(), status.GetGreen(), status.GetBlue())
	}
}

var Cmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(0),
	Use:   "get",
	Short: "Get LED status",
	Long:  "Get LED statuses from a remote blinky server.",
	Run:   getLEDs,
}
