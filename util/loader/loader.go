package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jucardi/infuse/util/log"
)

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

func LoadTemplate(filename string) (string, error) {
	log.Debug(" <-- loadtemplate entry")
	bs, err := ioutil.ReadFile(filename)
	return string(bs), err
}

func LoadJSON(jsonStr string) (map[string]interface{}, error) {
	log.Debug(" <-- LoadJSON entry")
	ret := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &ret)
	return ret, err
}
