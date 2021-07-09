package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
)

type config struct {
	debug     bool
	directory string
	output    string
}

type variableFile struct {
	cluster   string
	component string
	contents  map[string]interface{}
}

func getconfig() (*config, error) {
	cfg := &config{
	}

	pflag.StringVar(&cfg.directory, "directory", cfg.directory, "which directory to process")
	pflag.StringVar(&cfg.output, "output", cfg.output, "which directory to write to")
	pflag.BoolVar(&cfg.debug, "debug", cfg.debug, "enable debug output")
	pflag.Parse()

	if len(cfg.directory) == 0 {
		return nil, fmt.Errorf("--directory required")
	}
	if len(cfg.output) == 0 {
		return nil, fmt.Errorf("--output required")
	}

	return cfg, nil
}

func run() error {
	cfg, err := getconfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.debug {
		log.SetLevel(log.TraceLevel)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
		os.Exit(1)
	}
}
