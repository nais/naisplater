package templatetools

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Variables map[interface{}]interface{}

type CryptFunc func(source, key string) (result string, err error)

func CryptTransform(vars Variables, password string, fn CryptFunc, translate bool) error {
	for k, v := range vars {
		switch typed := v.(type) {
		case Variables:
			err := CryptTransform(typed, password, fn, translate)
			if err != nil {
				return fmt.Errorf("%s: %s", k, err)
			}
		case string:
			key, ok := k.(string)
			if !ok {
				return fmt.Errorf("non-string key '%v'", k)
			}
			if strings.HasSuffix(key, ".enc") {
				log.Debugf("Running crypt function on variable '%s'", key)
				plaintext, err := fn(typed, password)
				if err != nil {
					return fmt.Errorf("crypt error: %w", err)
				}
				if translate {
					vars[key[:len(key)-4]] = plaintext
					delete(vars, k)
				} else {
					vars[key] = plaintext
				}
			}
		}
	}

	return nil
}

func MergeMaps(dst, src Variables) error {
	for k, srcValue := range src {
		dstValue, ok := dst[k]
		if !ok {
			dst[k] = srcValue
			continue
		}
		dstMap, ok := dstValue.(Variables)
		if !ok {
			dst[k] = srcValue
			continue
		}
		srcMap, ok := srcValue.(Variables)
		if !ok {
			return fmt.Errorf("%s: trying to overwrite map variable with non-map type variable", k)
		}
		err := MergeMaps(dstMap, srcMap)
		if err != nil {
			return fmt.Errorf("%s: %s", k, err)
		}
	}

	return nil
}

func VariablesFromFiles(paths ...string) (Variables, error) {
	allVars := Variables{}
	vars := Variables{}

	for _, path := range paths {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: open file: %s", path, err)
		}

		err = yaml.Unmarshal(file, &vars)
		if err != nil {
			return nil, err
		}

		err = MergeMaps(allVars, vars)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
	}

	return allVars, nil
}
