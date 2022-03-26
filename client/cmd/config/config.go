package config

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vatine/blinky/client/pkg/protos"
)

var Server string
var LogLevel string

// Various post-init settings, common to one and all
func Setup() {
	log.SetLevel(log.InfoLevel)
	level, err := log.ParseLevel(LogLevel)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"string": LogLevel,
		}).Warn("failed to parse log level, continuing with info level")
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.WithFields(log.Fields{
		"server":    Server,
		"log level": LogLevel,
	}).Debug("setup complete")
}

func Connect() (protos.BlinkyClient, error) {
	conn, err := grpc.Dial(Server, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"server": Server,
		}).Error("failed to connect")
		return nil, err
	}

	return protos.NewBlinkyClient(conn), nil
}
