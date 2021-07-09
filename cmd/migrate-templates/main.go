package main

import (
	"fmt"
	"github.com/nais/naisplater/pkg/parser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	debug  bool
	input  string
	output string
}

func getconfig() (*config, error) {
	cfg := &config{
	}

	pflag.StringVar(&cfg.input, "input", cfg.input, "which directory to process")
	pflag.StringVar(&cfg.output, "output", cfg.output, "which directory to write to")
	pflag.BoolVar(&cfg.debug, "debug", cfg.debug, "enable debug output")
	pflag.Parse()

	if len(cfg.input) == 0 {
		return nil, fmt.Errorf("--input required")
	}
	if len(cfg.output) == 0 {
		return nil, fmt.Errorf("--output required")
	}

	return cfg, nil
}

func translate(inFile, outFile string) error {
	in, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	component := filepath.Base(inFile)
	ext := filepath.Ext(component)
	if strings.HasSuffix(component, ext) {
		component = component[:len(component)-len(ext)]
	}

	return parser.ReplaceVariables(in, out, "."+component)
}

func run() error {
	cfg, err := getconfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.debug {
		log.SetLevel(log.TraceLevel)
	}

	cfg.input = filepath.Clean(cfg.input)
	cfg.output, err = filepath.Abs(cfg.output)
	if err != nil {
		return err
	}

	walkFunc := func(path string, info fs.FileInfo, err error) error {
		dest := strings.Replace(path, cfg.input, cfg.output, 1)
		if info.IsDir() {
			log.Debugf("Create directory %s", dest)
			return os.MkdirAll(dest, 0755)
		}
		log.Debugf("Translating %s to %s", path, dest)
		err = translate(path, dest)
		if err != nil {
			log.Errorf("Error in %s: %s", path, err)
		}
		return nil
	}

	return filepath.Walk(cfg.input, walkFunc)
}

func main() {
	err := run()
	if err != nil {
		log.Errorf("fatal: %s", err)
		os.Exit(1)
	}
}
