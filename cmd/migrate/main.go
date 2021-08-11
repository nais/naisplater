package main

import (
	"fmt"
	"github.com/nais/naisplater/pkg/cryptutil"
	"github.com/nais/naisplater/pkg/templatetools"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"strings"

	yamlv2 "gopkg.in/yaml.v2"
)

type config struct {
	debug         bool
	directory     string
	output        string
	decryptionKey string
}

type variableFile struct {
	cluster   string
	component string
	contents  map[interface{}]interface{}
}

func getconfig() (*config, error) {
	cfg := &config{
	}

	pflag.StringVar(&cfg.directory, "directory", cfg.directory, "which directory to process")
	pflag.StringVar(&cfg.output, "output", cfg.output, "which directory to write to")
	pflag.StringVar(&cfg.decryptionKey, "decryption-key", cfg.decryptionKey, "decryption key for secrets")
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

func processDirectory(directory, key string) ([]*variableFile, error) {
	log.Debugf("found directory: %s", directory)

	dirEntry, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	results := make([]*variableFile, 0)

	for _, file := range dirEntry {
		path := filepath.Join(directory, file.Name())
		if file.IsDir() {
			sub, err := processDirectory(path, key)
			if err != nil {
				return nil, err
			}
			results = append(results, sub...)
		} else {
			component := filepath.Base(path)
			ext := filepath.Ext(component)
			if strings.HasSuffix(component, ext) {
				component = component[:len(component)-len(ext)]
			}
			result, err := processFile(filepath.Base(directory), component, path, key)
			if err != nil {
				return nil, err
			}
			log.Infof("Processed %s/%s with %d entries", result.cluster, result.component, len(result.contents))
			results = append(results, result)
		}
	}

	return results, nil
}

func processFile(cluster, component, path, key string) (*variableFile, error) {
	log.Debugf("found file: %s", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	result := &variableFile{
		cluster:   cluster,
		component: strings.Replace(component, "-", "_", -1),
	}

	dec := yamlv2.NewDecoder(file)
	err = dec.Decode(&result.contents)
	if err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	err = templatetools.CryptTransform(result.contents, key, cryptutil.ReEncrypt, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func clusters(results []*variableFile) []string {
	keys := make(map[string]interface{})
	output := make([]string, 0)
	for _, result := range results {
		keys[result.cluster] = new(interface{})
	}
	for cluster := range keys {
		output = append(output, cluster)
	}
	return output
}

func filter(cluster string, results []*variableFile) []*variableFile {
	filtered := make([]*variableFile, 0, len(results))
	for _, result := range results {
		if result.cluster == cluster {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func concat(results []*variableFile) map[string]interface{} {
	output := make(map[string]interface{})
	for _, result := range results {
		output[result.component] = result.contents
	}
	return output
}

func write(destination string, results []*variableFile) error {
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	defer file.Close()

	values := concat(results)

	enc := yamlv2.NewEncoder(file)
	return enc.Encode(values)
}

func run() error {
	cfg, err := getconfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.debug {
		log.SetLevel(log.TraceLevel)
	}

	results, err := processDirectory(cfg.directory, cfg.decryptionKey)
	if err != nil {
		return err
	}

	clusters := clusters(results)
	for _, cluster := range clusters {
		clusterResults := filter(cluster, results)
		destination := filepath.Join(cfg.output, cluster+".yaml")
		err = write(destination, clusterResults)
		if err != nil {
			return err
		}
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
