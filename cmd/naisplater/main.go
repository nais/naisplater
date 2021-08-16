package main

import (
	"bytes"
	"fmt"
	"github.com/nais/naisplater/pkg/cryptutil"
	"github.com/nais/naisplater/pkg/templatetools"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

type config struct {
	debug         bool
	encrypt       bool
	decrypt       string
	templates     string
	variables     string
	output        string
	cluster       string
	decryptionKey string
	addLabels     bool
	touchedAt     string
	validate      bool
}

func getconfig() (*config, error) {
	currentTime := time.Now()
	touchedAt := currentTime.Format("20060102T150405")

	cfg := &config{
		addLabels:     true,
		touchedAt:     touchedAt,
		decryptionKey: os.Getenv("NAISPLATER_DECRYPTION_KEY"),
	}

	pflag.StringVar(&cfg.templates, "templates", cfg.templates, "directory with templates")
	pflag.StringVar(&cfg.variables, "variables", cfg.variables, "directory with variables")
	pflag.StringVar(&cfg.output, "output", cfg.output, "which directory to write to")
	pflag.StringVar(&cfg.cluster, "cluster", cfg.cluster, "cluster for rendering templates and variables")
	pflag.StringVar(&cfg.decryptionKey, "decryption-key", cfg.decryptionKey, "key for decrypting variables ($NAISPLATER_DECRYPTION_KEY)")
	pflag.BoolVar(&cfg.debug, "debug", cfg.debug, "enable debug output")
	pflag.BoolVar(&cfg.addLabels, "add-labels", cfg.addLabels, "add 'nais.io/created-by' and 'nais.io/touched-at' labels")
	pflag.StringVar(&cfg.touchedAt, "touched-at", cfg.touchedAt, "use custom timestamp in 'nais.io/touched-at' label")
	pflag.BoolVar(&cfg.encrypt, "encrypt", cfg.encrypt, "in-place encrypt all plaintext values with 'key.enc' keys")
	pflag.StringVar(&cfg.decrypt, "decrypt", cfg.decrypt, "decrypt all ciphertext values with 'key.enc' keys in given file; output the whole file to STDOUT")
	pflag.BoolVar(&cfg.validate, "validate", cfg.validate, "render all templates for all clusters in-memory and check for syntax/runtime errors")
	pflag.Parse()

	if len(cfg.decrypt) == 0 && len(cfg.variables) == 0 {
		return nil, fmt.Errorf("--variables required")
	}
	if cfg.validate && (cfg.encrypt || len(cfg.decrypt) > 0) {
		return nil, fmt.Errorf("--validate cannot be used together with --encrypt or --decrypt")
	}
	if cfg.encrypt && len(cfg.decrypt) > 0 {
		return nil, fmt.Errorf("--encrypt and --decrypt are mutually exclusive")
	}
	if cfg.encrypt || len(cfg.decrypt) > 0 {
		if len(cfg.decryptionKey) == 0 {
			return nil, fmt.Errorf("--encrypt and --decrypt needs --decryption-key to work")
		}
		// return early for crypt-only operation
		return cfg, nil
	}
	if len(cfg.templates) == 0 {
		return nil, fmt.Errorf("--templates required")
	}
	if cfg.validate {
		// no --output or --cluster required for validation
		return cfg, nil
	}
	if len(cfg.cluster) == 0 {
		return nil, fmt.Errorf("--cluster required")
	}
	if len(cfg.output) == 0 {
		return nil, fmt.Errorf("--output required")
	}

	return cfg, nil
}

func render(inFile, outFile string, vars templatetools.Variables, cfg *config) error {
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

	log.Debugf("Rendering %s to %s", inFile, outFile)

	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, vars)
	if err != nil {
		return err
	}

	if !cfg.addLabels {
		_, err = io.Copy(out, buffer)
		return err
	}

	decoder := yaml.NewDecoder(buffer)
	encoder := yaml.NewEncoder(out)

	defer encoder.Close()

	bufbytes := buffer.Bytes()

	for {
		content := make(map[interface{}]interface{})
		err = decoder.Decode(&content)
		if err == io.EOF {
			encoder.Close()
			out.Close()
			return nil
		} else if err != nil {
			os.Stdout.Write([]byte("\n\n-----------------------\n\n"))
			os.Stdout.Write(bufbytes)
			return err
		}

		err = injectLabels(content, cfg.touchedAt)
		if err != nil {
			return err
		}

		err = encoder.Encode(content)
		if err != nil {
			return err
		}
	}

	return nil
}

func injectLabels(content map[interface{}]interface{}, touchedAt string) error {

	metadata, ok := content["metadata"].(map[interface{}]interface{})
	if !ok {
		return nil
	}

	labels, ok := metadata["labels"].(map[interface{}]interface{})
	if !ok {
		metadata["labels"] = make(map[interface{}]interface{})
		labels = metadata["labels"].(map[interface{}]interface{})
	}

	labels["nais.io/created-by"] = "nais-yaml"
	labels["nais.io/touched-at"] = touchedAt

	return nil
}

func variablefilename(cluster string) string {
	if len(cluster) == 0 {
		return "vars.yaml"
	}
	return cluster + ".yaml"
}

func merge(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

func directoryTemplates(directory string) (map[string]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		basename := entry.Name()
		files[basename] = filepath.Join(directory, basename)
	}

	return files, nil
}

func allClusters(cfg *config) ([]string, error) {
	dirEntry, err := os.ReadDir(cfg.variables)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	clusters := make([]string, 0)
	for _, file := range dirEntry {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".yaml") {
			clusters = append(clusters, file.Name()[:len(file.Name())-5])
		}
	}

	return clusters, nil
}

func encrypt(cfg *config) error {
	dirEntry, err := os.ReadDir(cfg.variables)
	if err != nil {
		return fmt.Errorf("read directory: %w", err)
	}

	for _, file := range dirEntry {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(cfg.variables, file.Name())
		vars, err := templatetools.VariablesFromFiles(path)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		err = templatetools.CryptTransform(vars, cfg.decryptionKey, cryptutil.EncryptIfPlaintext, false)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		tmpfile, err := os.CreateTemp(os.TempDir(), "naisplater")
		if err != nil {
			return err
		}
		err = yaml.NewEncoder(tmpfile).Encode(vars)
		if err != nil {
			return err
		}
		err = tmpfile.Close()
		if err != nil {
			return err
		}
		err = os.Rename(tmpfile.Name(), path)
		if err != nil {
			return err
		}
	}

	return nil
}

func decrypt(cfg *config) error {
	vars, err := templatetools.VariablesFromFiles(cfg.decrypt)
	if err != nil {
		return err
	}

	err = templatetools.CryptTransform(vars, cfg.decryptionKey, cryptutil.DecryptWithPassword, false)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(os.Stdout).Encode(vars)
}

func run(cfg *config) error {
	errors := 0

	globals := filepath.Join(cfg.variables, variablefilename(""))
	locals := filepath.Join(cfg.variables, variablefilename(cfg.cluster))

	log.Debugf("Using global variables from %s", globals)
	log.Debugf("Using cluster-override variables from %s", locals)

	vars, err := templatetools.VariablesFromFiles(globals, locals)
	if err != nil {
		return err
	}

	log.Debugf("Decrypting variables")
	err = templatetools.CryptTransform(vars, cfg.decryptionKey, cryptutil.DecryptWithPassword, true)
	if err != nil {
		if len(cfg.decryptionKey) == 0 {
			log.Errorf("decrypt variable: %s", err)
			log.Warnf("Decryption key is missing; skipping all variable decryption")
			errors++
		} else {
			return err
		}
	}

	log.Debugf("Using templates from %s", cfg.templates)

	templates, err := directoryTemplates(cfg.templates)
	if err != nil {
		return err
	}

	clusterTemplates := filepath.Join(cfg.templates, cfg.cluster)
	overrides, err := directoryTemplates(clusterTemplates)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("No cluster-specific template directory for '%s'", cfg.cluster)
		} else {
			return err
		}
	} else {
		log.Debugf("Using cluster-override templates from %s", clusterTemplates)
	}

	merge(templates, overrides)

	filenames := make([]string, 0, len(templates))
	for k := range templates {
		filenames = append(filenames, k)
	}
	sort.Strings(filenames)

	if !cfg.validate {
		log.Debugf("Using output directory %s", cfg.output)
		err = os.MkdirAll(cfg.output, 0755)
		if err != nil {
			return err
		}
	}

	for _, filename := range filenames {
		path := templates[filename]
		output := filepath.Join(cfg.output, filename)
		if cfg.validate {
			output = "/dev/null"
		}
		err = render(path, output, vars, cfg)
		if err != nil {
			errors++
			log.Errorf("Render %s: %s", path, err)
		} else {
			log.Debugf("Rendered %s", output)
		}
	}

	if errors > 0 {
		return fmt.Errorf("encountered %d errors; see log", errors)
	}

	return nil
}

func validate(cfg *config) error {
	clusters, err := allClusters(cfg)
	if err != nil {
		return err
	}

	errors := 0
	for _, cluster := range clusters {
		log.Infof("Running validation for cluster '%s'", cluster)

		cfg.cluster = cluster
		err = run(cfg)

		if err != nil {
			log.Errorf("Validation failed for cluster '%s': %s", cluster, err)
			errors++
		}
	}

	if errors > 0 {
		return fmt.Errorf("%d clusters failed validation", errors)
	}

	log.Infof("All clusters rendered successfully")

	return nil
}

func runner() error {
	cfg, err := getconfig()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.debug {
		log.SetLevel(log.TraceLevel)
	}

	if cfg.encrypt {
		return encrypt(cfg)
	}

	if len(cfg.decrypt) > 0 {
		return decrypt(cfg)
	}

	if cfg.validate {
		return validate(cfg)
	}

	return run(cfg)
}

func main() {
	err := runner()
	if err != nil {
		log.Errorf("fatal: %s", err)
		os.Exit(1)
	}
}
