package spectrum

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vatine/blinky/client/cmd/config"
	"github.com/vatine/blinky/client/pkg/protos"
)

var Cmd cobra.Command
var sleep time.Duration
var brightness int

type color struct {
	r, g, b float64
}

// Generate a sweep starting at one color, finishing at another,
// having (in total) a certain number of steps.
func interpolate(start, end color, steps int) []color {
	s := float64(steps - 1)

	dr := (end.r - start.r) / s
	dg := (end.g - start.g) / s
	db := (end.b - start.b) / s

	rv := []color{start}
	for ds := 1.0; ds <= s; ds += 1.0 {
		next := color{r: start.r + dr*ds, g: start.g + dg*ds, b: start.b + db*ds}
		rv = append(rv, next)
	}

	log.WithFields(log.Fields{
		"rv-len": len(rv),
		"steps":  steps,
		"rv[-1]": rv[len(rv)-1],
		"end":    end,
		"rv":     rv,
	}).Debug("interpolation done")

	return rv
}

// For now, assume a 32-LED blinky
func draw(client protos.BlinkyClient, max int32) {
	var cols []color

	r_o := interpolate(color{r: 1.0}, color{r: 1.0, g: 0.5}, 6)
	cols = append(cols, r_o...)

	o_y := interpolate(color{r: 1.0, g: 0.5}, color{r: 1.0, g: 1.0}, 6)
	cols = append(cols, o_y...)

	y_g := interpolate(color{r: 1.0, g: 1.0}, color{g: 1.0}, 7)
	cols = append(cols, y_g...)

	g_b := interpolate(color{g: 1.0}, color{b: 1.0}, 7)
	cols = append(cols, g_b...)

	b_i := interpolate(color{b: 1.0}, color{r: 0.294, b: 0.509}, 4)
	cols = append(cols, b_i...)

	i_v := interpolate(color{r: 0.294, b: 0.509}, color{r: 0.58, b: 0.827}, 2)
	cols = append(cols, i_v...)
	log.WithFields(log.Fields{
		"len": len(cols),
	}).Debug("bow")

	fmax := float64(max)
	for led, col := range cols {

		r := int32(fmax * col.r)
		g := int32(fmax * col.g)
		b := int32(fmax * col.b)

		req := protos.SetLEDRequest{LEDs: []int32{int32(led)}, Red: r, Green: g, Blue: b}
		client.SetLEDs(context.Background(), &req)
	}
}

func spectrum(cmd *cobra.Command, args []string) {
	config.Setup()

	client, err := config.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("spectrum - failed to connect")
		return
	}

	draw(client, int32(brightness))
}

func init() {
	Cmd.Run = spectrum
	Cmd.Use = "spectrum"
	Cmd.Short = "make spectrum"
	Cmd.Long = "Make an approximate spectrum"
	Cmd.Flags().IntVar(&brightness, "brightness", 70, "Maximum intensity on each R/G/B channel.")
}
