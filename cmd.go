package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/jessevdk/go-flags"
)

// Command-line flag options
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose logging."`
	Number  int    `long:"num" description:"Number of concurrent hammers."`
}

var logLevels = []log.Level{
	log.Warning,
	log.Info,
	log.Debug,
}

func main() {
	// Setup cli parser
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.Usage = "[OPTIONS] hostname[:port]"

	args, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
		return
	} else if len(args) < 1 {
		fmt.Println("Invalid usage: Missing hostname.")
		return
	}

	host := args[0]
	if !strings.Contains(host, ":") {
		host += ":22"
	}

	// Set verbosity of logger
	numVerbose := len(options.Verbose)
	if numVerbose > len(logLevels) {
		numVerbose = len(logLevels)
	}

	logLevel := logLevels[numVerbose]
	logger = golog.New(os.Stderr, logLevel)

	// Construct interrupt handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	logger.Infof("Hammering: %s", host)
	h := NewHammer(host, options.Number)

	err = h.Start()
	if err != nil {
		logger.Errorf("Failed to start: %s", err)
		return
	}

	<-sig // Wait for ^C signal
	logger.Warningf("Interrupt signal detected, shutting down.")
	h.Stop()
}
