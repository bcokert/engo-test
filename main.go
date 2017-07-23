package main

import (
	"os"

	"engo.io/engo"
	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/scenes/owlclicker"
)

func main() {
	options := engo.RunOptions{
		Title:  "Owl Game",
		Width:  800, // pixels
		Height: 600, // pixels
	}

	var logLevel logging.LOGLEVEL
	switch os.Getenv("ENGOLOGLEVEL") {
	case "DEBUG":
		logLevel = logging.DEBUG
	default:
		logLevel = logging.INFO
	}
	logger := logging.NewDefaultLogger(logLevel)

	engo.Run(options, &owlclicker.Scene{Log: logger})
}
