package templatetools

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Variables map[interface{}]interface{}

func Decrypt(vars Variables) error {
	for k, v := range vars {
		switch typed := v.(type) {
		case Variables:
			err := Decrypt(typed)
			if err != nil {
				return fmt.Errorf("%s: %s", k, err)
			}
		case string:
			key, ok := k.(string)
			if !ok {
				return fmt.Errorf("non-string key '%v'", k)
			}
			if strings.HasSuffix(key, ".enc") {
				log.Debugf("Decrypting variable '%s'", key)
				vars[key[:len(key)-4]] = v
				delete(vars, k)
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
