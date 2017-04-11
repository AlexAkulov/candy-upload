package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
	"github.com/spf13/pflag"
)

var (
	version   = "unknown"
	goVersion = "unknown"
	buildDate = "unknown"

	config *Config
	log    *logging.Logger
)

func main() {

	versionFlag := pflag.BoolP("version", "v", false, "Print version and exit")
	configPath := pflag.StringP("config", "c", "config.yml", "Path to config file")
	helpFlag := pflag.BoolP("help", "h", false, "Print this message and exit")

	pflag.Parse()

	if *helpFlag {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println("version: ", version)
		fmt.Println("Goland version: ", goVersion)
		fmt.Println("Build Date: ", buildDate)
		os.Exit(0)
	}
	var err error
	if config, err = loadConfig(*configPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if log, err = newLog(config.LogFile, config.LogLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := &Server{
		config: config,
		log:    log,
	}

	if err := s.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	log.Infof("Candy Upload started (%s)", version)

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Received signal: ", <-shutdown)
	s.Stop()
	log.Infof("Candy Upload stopped (%s)", version)
}
