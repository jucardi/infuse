package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jucardi/infuse/util/log"
	"gopkg.in/yaml.v2"
)

// LoadTemplates loads multiple template files and returns a map of filename,value
func LoadTemplates(searchArg string) (map[string]string, error) {
	log.Debug(" <-- loadtemplates entry")
	ret := map[string]string{}
	matches, err := filepath.Glob(searchArg)

	if err != nil {
		return nil, err
	}

	for _, f := range matches {
		inf, err := os.Stat(f)

		if err != nil {
			return nil, fmt.Errorf("unable to read '%s'", f)
		} else if inf.IsDir() {
			continue
		}

		if str, err := LoadTemplate(f); err != nil {
			return nil, fmt.Errorf("failed to load '%s'", f)
		} else {
			ret[inf.Name()] = str
		}
	}

	return ret, nil
}

// LoadTemplate loads a file template
func LoadTemplate(filename string) (string, error) {
	log.Debug(" <-- loadtemplate entry")
	bs, err := ioutil.ReadFile(filename)
	return string(bs), err
}

// LoadMarshalled attempts to unmarshall the byte data provided to a map[string]interface{}, first will try to unmarshall as JSON and if fails will attempt to unmarshall as YAML
func LoadMarshalled(data []byte) (map[string]interface{}, error) {
	ret, jsonErr := LoadJSON(data)
	if jsonErr == nil {
		return ret, nil
	}
	ret, yamlErr := LoadYAML(data)
	if yamlErr == nil {
		return ret, nil
	}
	return nil, fmt.Errorf("failed to unmarshall as JSON and YAML.\n- JSON error: %+v\n- YAML error: %+v", jsonErr, yamlErr)
}

// LoadJSON attempts to unmarshall the byte data provided representing a JSON object to a map[string]interface{}
func LoadJSON(data []byte) (ret map[string]interface{}, err error) {
	ret = map[string]interface{}{}
	err = json.Unmarshal(data, &ret)
	return
}

// LoadYAML attempts to unmarshall the byte data provided representing a YAML object to a map[string]interface{}
func LoadYAML(data []byte) (ret map[string]interface{}, err error) {
	ret = map[string]interface{}{}
	err = yaml.Unmarshal(data, &ret)
	return
}
