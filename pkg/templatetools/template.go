package templatetools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Variables map[string]interface{}

func VariablesFromFiles(paths ...string) (Variables, error) {
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
	}

	return vars, nil
}
