package fade

import (
	"context"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/pkg/protos"
)

func fadeLEDs(cmd *cobra.Command, args []string) {
	config.Setup()

	log.WithFields(log.Fields{
		"leds": config.LEDs,
	}).Debug("fade called")

	req := protos.GetLEDRequest{LEDs: config.LEDs}
	client, err := config.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get client")
		return
	}

	resp, err := client.GetLEDs(context.Background(), &req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error()
	}

	buildFadeRequest(client, resp)
}

func buildFadeRequest(client protos.BlinkyClient, resp *protos.GetLEDResponse) {
	statuses := resp.GetStatus()
	ctx := context.Background()

	low, high := getRange(statuses)
	log.WithFields(log.Fields{
		"low":  low,
		"high": high,
	}).Debug("status")

	if high == low {
		log.Info("Not much sense in fading without a span")
		return
	}

	startRed := statuses[low-low].GetRed()
	startGreen := statuses[low-low].GetGreen()
	startBlue := statuses[low-low].GetBlue()
	endRed := statuses[high-low].GetRed()
	endGreen := statuses[high-low].GetGreen()
	endBlue := statuses[high-low].GetBlue()

	floatCount := float64(high - low)
	deltaRed := float64(endRed-startRed) / floatCount
	deltaGreen := float64(endGreen-startGreen) / floatCount
	deltaBlue := float64(endBlue-startBlue) / floatCount

	for led := low + 1; led < high; led++ {
		scaler := float64(led - low)
		newRed := float64(startRed) + scaler*deltaRed
		newGreen := float64(startGreen) + scaler*deltaGreen
		newBlue := float64(startBlue) + scaler*deltaBlue

		req := protos.SetLEDRequest{
			LEDs:  []int32{led},
			Red:   int32(math.Round(newRed)),
			Green: int32(math.Round(newGreen)),
			Blue:  int32(math.Round(newBlue)),
		}
		log.WithFields(log.Fields{
			"req": req,
		}).Debug("fade set request")

		client.SetLEDs(ctx, &req)
	}

}

func getRange(status []*protos.LEDStatus) (int32, int32) {
	var low, high int32
	low = 255
	high = 0

	for _, s := range status {
		led := s.GetLED()
		if led < low {
			low = led
		}
		if led > high {
			high = led
		}
	}

	return low, high
}

var Cmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(0),
	Use:   "fade",
	Short: "Fade LEDs",
	Long:  "Fade LEDs on a remote blinky server",
	Run:   fadeLEDs,
}
