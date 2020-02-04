package gotmpl

import (
	"encoding/json"
	"fmt"
	"github.com/jucardi/go-osx/paths"
	"github.com/jucardi/infuse/templates/helpers"
	"github.com/jucardi/infuse/util/log"
	"gopkg.in/jucardi/go-streams.v1/streams"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"
)

var instance *helperContext

func getHelpers() *helperContext {
	if instance == nil {
		instance = &helperContext{
			IHelpersManager: helpers.New(),
		}
		instance.init()
	}
	return instance
}

type helperContext struct {
	helpers.IHelpersManager
	*template.Template
}

func (h *helperContext) setTemplate(tmpl *template.Template) {
	h.Template = tmpl
}

func (h *helperContext) toMap() template.FuncMap {
	ret := template.FuncMap{}
	for _, v := range h.Get() {
		ret[v.Name] = v.Function
	}
	return ret
}

func (h *helperContext) init() {
	helpers.RegisterCommon(h)
	_ = h.Register("default", h.defaultFn, "The first argument should be a default value, and the second argument is a value that will be evaluated. If arg2 is a zero value, returns arg1, otherwise returns arg2")
	_ = h.Register("map", h.mapFn, "Creates a new map[string]interface{}, the provided arguments should be key, value, key, value...")
	_ = h.Register("dict", h.mapFn, "Creates a new map[string]interface{}, the provided arguments should be key, value, key, value...")
	_ = h.Register("include", h.includeFile, "Includes a template file as an internal template reference by the provided name")
	_ = h.Register("set", h.setFn, "Allows to set a value to a map[string]interface{}")
	_ = h.Register("append", h.append, "Appends a value into an existing array")
	_ = h.Register("iterate", h.iterate, "Creates an iteration array of the provided length, so it can be used as {{ range $val := iterate N }} where N is the length of the iteration. Created due to the lack of `for` loops.")
	_ = h.Register("loadJson", h.loadJson, "Unmarshals a JSON string into a map[string]interface{}")
}

func (h *helperContext) defaultFn(val ...interface{}) interface{} {
	for i := len(val) - 1; i > 0; i-- {
		x := val[i]
		v := reflect.ValueOf(x)

		switch v.Kind() {
		case reflect.String:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			if v.Len() > 0 {
				return x
			}
		default:
			if v.IsValid() {
				return x
			}
		}

	}
	return val[0]
}

func (h *helperContext) mapFn(args ...interface{}) map[string]interface{} {
	if len(args)%2 != 0 {
		log.Panicf("Error in 'map' directive. The number of keys do not match the number of values")
	}

	ret := map[string]interface{}{}

	for i := 0; i < len(args); i = i + 2 {
		ret[fmt.Sprintf("%v", args[i])] = args[i+1]
	}

	return ret
}

func (h *helperContext) includeFile(name, file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting work directory %s", err.Error())
	}
	templateData, err := ioutil.ReadFile(paths.Combine(wd, file))
	if err != nil {
		return "", fmt.Errorf("error including template file %s, %s", file, err.Error())
	}
	tmpl, err := h.Template.New(name).Parse(string(templateData))
	if err != nil {
		return "", fmt.Errorf("error parsing template file %s, %s", file, err.Error())
	}
	h.Template = tmpl
	return "", nil
}

func (h *helperContext) setFn(obj interface{}, key string, value interface{}) string {
	m := obj.(map[string]interface{})
	m[key] = value
	return ""
}

func (h *helperContext) append(array interface{}, values ...interface{}) interface{} {
	vals := streams.From(values).
		Map(func(i interface{}) interface{} {
			return reflect.ValueOf(i)
		}).
		ToArray().([]reflect.Value)
	return reflect.Append(reflect.ValueOf(array), vals...).Interface()
}

func (h *helperContext) iterate(count int) []int {
	var array []int
	for i := 0; i < count; i++ {
		array = append(array, i)
	}
	return array
}

func (h *helperContext) loadJson(str string) map[string]interface{} {
	ret := map[string]interface{}{}
	if err := json.Unmarshal([]byte(str), &ret); err != nil {
		panic(err.Error())
	}
	return ret
}
