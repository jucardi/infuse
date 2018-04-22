package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jucardi/gotmpl/util/log"
)

func loadTemplates(searchArg string) (map[string]string, error) {
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

		if str, err := loadTemplate(f); err != nil {
			return nil, fmt.Errorf("failed to load '%s'", f)
		} else {
			ret[inf.Name()] = str
		}
	}

	return ret, nil
}

func loadTemplate(filename string) (string, error) {
	log.Debug(" <-- loadtemplate entry")
	bs, err := ioutil.ReadFile(filename)
	return string(bs), err
}

func loadJSON(jsonStr string) (map[string]interface{}, error) {
	log.Debug(" <-- loadJSON entry")
	ret := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &ret)
	return ret, err
}
