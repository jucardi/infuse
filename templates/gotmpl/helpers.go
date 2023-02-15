package gotmpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"text/template"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/infuse/templates/helpers"
	"github.com/jucardi/infuse/util/log"
	"github.com/jucardi/infuse/util/maps"
	"github.com/jucardi/infuse/util/reflectx"
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
	loaded map[string]*template.Template
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
	_ = h.Register("includeAsString", h.includeTemplate, "Includes a provided template string as an internal template reference by the provided name. Eg: {{ include [name] [contents] }}")
	_ = h.Register("set", h.setFn, "Allows to set a value to a map[string]interface{} or map[interface{}]interface{}")
	_ = h.Register("append", h.append, "Appends a value into an existing array")
	_ = h.Register("iterate", h.iterate, "Creates an iteration array of the provided length, so it can be used as {{ range $val := iterate N }} where N is the length of the iteration. Created due to the lack of `for` loops.")
	_ = h.Register("loadJson", h.loadJson, "Unmarshals a JSON string into a map[string]interface{}")
	_ = h.Register("mapSet", h.mapSetFn, `Allows to set a value using an XPATH representation of the key. Accepts an optional argument to indicate if the parents should be created if they don't exist'. E.g: {{mapSet $map ".some.key.path" $value $makeEmpty }}`)
	_ = h.Register("mapGet", h.mapGetFn, `Allows to get a value from a map using an XPATH representation of the key. Accepts optional argument for a default value to return if the value is not found". E.g: {{mapGet $map ".some.key.path" $someDefaultValue }}`)
	_ = h.Register("mapContains", h.mapContainsFn, `Indicates whether a value at the provided XPATH representation of the key exists in the provided map`)
	_ = h.Register("mapConvert", h.mapConvertFn, `Ensures the provided map is map[string]interface{}. Useful when loading values from a YAML where the deserialization is map[interface{}]interface{}`)
	_ = h.Register("invoke", h.invoke, `Similar to {{ template [name] [data] }}, invokes a name by the given name with the given data. The difference with 'template' is that 'invoke' can be used with a string value as the name instead of a hardcoded string`)
	_ = h.Register("parse", h.parse, `Attempts to parse the provided template contents using the provided data object and returns the parsed value. Usage {{ parse [obj] [template contents] }}`)
	_ = h.Register("parseXpath", h.parseXpath, `Attempts to parse a value inside a data object as a template and returns the parsed value using the same entry object to parse the template. Usage {{ parse [obj] [xpath to value with template] }}`)
	_ = h.Register("in", h.in, `Indicates whether a value is contained in an array. Usage:  {{ in [array] [value] }}`)
}

func (h *helperContext) mapSetFn(obj interface{}, key string, value interface{}, makeEmpty ...bool) string {
	var inMap map[string]interface{}

	switch m := obj.(type) {
	case map[string]interface{}:
		inMap = m
	case map[interface{}]interface{}:
		if converted, err := maps.ConvertMap(obj); err != nil {
			panic(fmt.Sprintf("failed to convert map[interface{}]interface{} to map[string]interface{}, %s", err.Error()))
		} else {
			inMap = converted
		}
	}

	if inMap == nil {
		panic(fmt.Sprintf("type not supported for map operations %T", obj))
	}

	if err := maps.SetValue(inMap, key, value, len(makeEmpty) > 0 && makeEmpty[0]); err != nil {
		panic(fmt.Sprintf("failed to set value to map using key '%s'  >  %s", key, err.Error()))
	}

	return ""
}

func (h *helperContext) mapGetFn(obj interface{}, key string, defaultValue ...interface{}) interface{} {
	var (
		inMap map[string]interface{}
		ret   interface{}
	)

	switch m := obj.(type) {
	case map[string]interface{}:
		inMap = m
	case map[interface{}]interface{}:
		if converted, err := maps.ConvertMap(obj); err != nil {
			panic(fmt.Sprintf("failed to convert map[interface{}]interface{} to map[string]interface{}, %s", err.Error()))
		} else {
			inMap = converted
		}
	}

	if len(defaultValue) > 0 {
		ret = defaultValue[0]
	}

	if inMap == nil {
		return ret
	}

	return maps.GetOrDefault(inMap, key, ret)
}

func (h *helperContext) mapContainsFn(obj interface{}, key string) bool {
	var inMap map[string]interface{}

	switch m := obj.(type) {
	case map[string]interface{}:
		inMap = m
	case map[interface{}]interface{}:
		if converted, err := maps.ConvertMap(obj); err != nil {
			panic(fmt.Sprintf("failed to convert map[interface{}]interface{} to map[string]interface{}, %s", err.Error()))
		} else {
			inMap = converted
		}
	}
	return maps.Contains(inMap, key)
}

func (h *helperContext) mapConvertFn(obj interface{}) map[string]interface{} {
	ret, err := maps.ConvertMap(obj)
	if err != nil {
		panic(fmt.Sprintf("failed to convert to map[string]interface{}, %s", err.Error()))
	}
	return ret
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
	templateData, err := ioutil.ReadFile(file)
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

func (h *helperContext) includeTemplate(name, contents string) (string, error) {
	tmpl, err := h.Template.New(name).Parse(contents)
	if err != nil {
		return "", fmt.Errorf("error parsing template by name %s, %s", name, err.Error())
	}
	h.Template = tmpl
	return "", nil
}

func (h *helperContext) setFn(obj interface{}, key string, value interface{}) string {
	switch m := obj.(type) {
	case map[string]interface{}:
		m[key] = value
	case map[interface{}]interface{}:
		m[key] = value
	}

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

func (h *helperContext) invoke(name string, data interface{}) (string, error) {
	tmpl := h.Template.Lookup(name)
	if tmpl == nil {
		return "", fmt.Errorf("failed to invoke template '%s', not found", name)
	}
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	return buf.String(), err
}

func (h *helperContext) parse(data interface{}, templateData string, failOnEmptyResult ...bool) (string, error) {
	if templateData == "" {
		if len(failOnEmptyResult) > 0 && failOnEmptyResult[0] {
			return "", errors.New("template produced empty result")
		}
		return "", nil
	}
	name := "__internal/parse/" + templateData
	tmpl := h.Template.Lookup(name)
	if tmpl == nil {
		if _, err := h.includeTemplate(name, templateData); err != nil {
			return "", err
		}
		tmpl = h.Template.Lookup(name)
	}
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	ret := buf.String()
	if len(failOnEmptyResult) > 0 && failOnEmptyResult[0] && ret == "" && err == nil {
		return "", errors.New("template produced empty result")
	}
	return ret, err
}

func (h *helperContext) parseXpath(data interface{}, xpath string, failOnEmptyResult ...bool) (string, error) {
	templateData, ok := h.mapGetFn(data, xpath).(string)
	if !ok {
		return "", fmt.Errorf("failed to obtain template data, the provided object does not contain a string at the provided key '%s'", xpath)
	}
	return h.parse(data, templateData, failOnEmptyResult...)
}

func (h *helperContext) in(array interface{}, value interface{}) bool {
	arrVal := reflect.ValueOf(array)
	if kind := arrVal.Kind(); kind != reflect.Slice && kind != reflect.Array {
		panic("attempting to use 'in' with a non-array type")
	}
	if reflectx.IsNil(arrVal) || arrVal.Len() == 0 {
		return false
	}
	return streams.From(array).Contains(value)
}
