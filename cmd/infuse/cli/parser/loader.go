package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

var zeroVal = reflect.Value{}

type Data map[string]interface{}

func (d Data) LoadContents(contents []byte, file string) error {
	var err error
	split := strings.Split(file, ".")
	ext := strings.ToLower(split[len(split)-1])
	val := map[string]interface{}{}

	switch ext {
	case "json":
		err = json.Unmarshal(contents, &val)
	case "yml":
		fallthrough
	case "yaml":
		err = yaml.Unmarshal(contents, val)
	default:
		return fmt.Errorf("unknown file type %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to unmarshal %s, %s", file, err.Error())
	}

	merge(d, val)
	return nil
}

func (d Data) LoadFile(file string) error {
	_, filename := filepath.Split(file)
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s, %s", filename, err.Error())
	}
	return d.LoadContents(contents, filename)
}

func (d Data) LoadURL(url string) error {
	if resp, err := http.Get(url); err != nil {
		return fmt.Errorf("error occurred while fetching from URL, %+v", err)
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unsuccessful response code (%d)", resp.StatusCode)
	} else {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		split := strings.Split(url, "/")
		split = strings.Split(split[len(split)-1], "?")
		return d.LoadContents(buf.Bytes(), split[0])
	}
}

func (d Data) ToMap() map[string]interface{} {
	return map[string]interface{}(d)
}

func merge(dest interface{}, source interface{}) {
	dVal, ok := dest.(reflect.Value)
	if !ok {
		dVal = reflect.ValueOf(dest)
	}
	sVal, ok := source.(reflect.Value)
	if !ok {
		sVal = reflect.ValueOf(source)
	}

	if sVal.Type().Kind() != reflect.Map || dVal.Type().Kind() != reflect.Map {
		return
	}

	for _, k := range sVal.MapKeys() {
		val := reflect.ValueOf(sVal.MapIndex(k).Interface())

		if val.Kind() == reflect.Map || val.Kind() == reflect.Struct {
			isZero, target := true, reflect.Value{}
			if tVal := dVal.MapIndex(k); tVal.IsValid() {
				target = reflect.ValueOf(dVal.MapIndex(k).Interface())
				isZero = target == zeroVal
			}

			if isZero {
				dVal.SetMapIndex(k, val)
			} else {
				merge(target, val)
			}
		} else {
			dVal.SetMapIndex(k, val)
		}
	}
}
