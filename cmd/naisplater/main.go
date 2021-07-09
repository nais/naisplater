package main

import (
	"fmt"
	"github.com/nais/naisplater/pkg/templatetools"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"text/template"
)

type config struct {
	debug     bool
	templates string
	variables string
	output    string
}

func getconfig() (*config, error) {
	cfg := &config{
	}

	pflag.StringVar(&cfg.templates, "templates", cfg.templates, "directory with templates")
	pflag.StringVar(&cfg.variables, "variables", cfg.variables, "directory with variables")
	pflag.StringVar(&cfg.output, "output", cfg.output, "which directory to write to")
	pflag.BoolVar(&cfg.debug, "debug", cfg.debug, "enable debug output")
	pflag.Parse()

	if len(cfg.templates) == 0 {
		return nil, fmt.Errorf("--templates required")
	}
	if len(cfg.variables) == 0 {
		return nil, fmt.Errorf("--variables required")
	}
	if len(cfg.output) == 0 {
		return nil, fmt.Errorf("--output required")
	}

	return cfg, nil
}

func render(inFile, outFile string, vars templatetools.Variables) error {
	out, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	tpl, err := template.ParseFiles(inFile)
	if err != nil {
		return err
	}

	// Nice API. Fail on undefined template variables.
	tpl.Option("missingkey=error")

	return tpl.Execute(out, vars)
}

func run() error {
	cfg, err := getconfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.debug {
		log.SetLevel(log.TraceLevel)
	}

	vars, err := templatetools.VariablesFromFiles(cfg.variables)
	if err != nil {
		return err
	}

	return render(cfg.templates, cfg.output, vars)
}

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
		os.Exit(1)
	}
}
