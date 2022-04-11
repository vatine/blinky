package runner

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/pkg/protos"
)

var Cmd cobra.Command
var sweeps int
var sleep time.Duration
var targetRed, targetGreen, targetBlue int

func half(red, green, blue int32, source *protos.LEDStatus) (int32, int32, int32) {
	r := (red + source.Red) / 2
	g := (green + source.Green) / 2
	b := (blue + source.Blue) / 2

	return r, g, b
}

func runner(cmd *cobra.Command, args []string) {
	config.Setup()

	client, err := config.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("runner - failed to connect")
		return
	}

	leds, err := client.GetLEDs(context.Background(), &protos.GetLEDRequest{})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failure getting LEDs")
	}
	red, _ := strconv.Atoi(args[0])
	green, _ := strconv.Atoi(args[1])
	blue, _ := strconv.Atoi(args[2])
	tRed := int32(targetRed)
	tGreen := int32(targetGreen)
	tBlue := int32(targetBlue)

	ledCount := len(leds.Status)
	for ix := 0; ix < sweeps; ix++ {
		for led := 0; led < ledCount; led++ {
			req := protos.SetLEDRequest{
				Red:   int32(red),
				Green: int32(green),
				Blue:  int32(blue),
				LEDs:  []int32{int32(led)},
			}
			client.SetLEDs(context.Background(), &req)
			if led > 0 {
				r, g, b := half(tRed, tGreen, tBlue, leds.Status[led-1])
				req := protos.SetLEDRequest{
					Red:   r,
					Green: g,
					Blue:  b,
					LEDs:  []int32{int32(led - 1)},
				}
				log.WithFields(log.Fields{
					"red":   r,
					"green": g,
					"blue":  b,
					"led":   led - 1,
				}).Debug()
				client.SetLEDs(context.Background(), &req)
			}
			time.Sleep(sleep)
		}
		for led := ledCount - 1; led >= 0; led-- {
			req := protos.SetLEDRequest{
				Red:   int32(red),
				Green: int32(green),
				Blue:  int32(blue),
				LEDs:  []int32{int32(led)},
			}
			client.SetLEDs(context.Background(), &req)
			req = protos.SetLEDRequest{
				Red:   tRed,
				Green: tGreen,
				Blue:  tBlue,
				LEDs:  []int32{int32(led + 1)},
			}
			client.SetLEDs(context.Background(), &req)
			time.Sleep(sleep)
		}
		req := protos.SetLEDRequest{
			Red:   0,
			Green: 0,
			Blue:  0,
			LEDs:  []int32{0},
		}
		client.SetLEDs(context.Background(), &req)
	}
}

func init() {
	Cmd.Args = cobra.MinimumNArgs(3)
	Cmd.Run = runner
	Cmd.Use = "run"
	Cmd.Short = "light-runner"
	Cmd.Long = "Light-runner"
	Cmd.Flags().IntVar(&sweeps, "sweeps", 1, "Number of back-and-forth sweeps.")
	Cmd.Flags().DurationVar(&sleep, "pause", time.Second/25, "Time to sleep between steps in a sweep.")
	Cmd.Flags().IntVar(&targetRed, "red", 0, "Post-sweep red.")
	Cmd.Flags().IntVar(&targetGreen, "green", 0, "Post-sweep green.")
	Cmd.Flags().IntVar(&targetBlue, "blue", 0, "Post-sweep blue.")
}
