package gotmpl

import (
	"fmt"
	"github.com/jucardi/infuse/templates/helpers"
	"github.com/jucardi/infuse/util/log"
	"reflect"
	"text/template"
)

var instance helpers.IHelpersManager

// Helpers returns the singleton helpers instance used for Go templates
func Helpers() helpers.IHelpersManager {
	if instance == nil {
		instance = helpers.New()
	}
	return instance
}

func getHelpers() template.FuncMap {
	ret := template.FuncMap{}
	for _, v := range Helpers().Get() {
		ret[v.Name] = v.Function
	}
	return ret
}

func init() {
	helpers.RegisterCommon(Helpers())
	_ = Helpers().Register("default", defaultFn, "The first argument should be a default value, and the second argument is a value that will be evaluated. If arg2 is a zero value, returns arg1, otherwise returns arg2")
	_ = Helpers().Register("map", mapFn, "Creates a new map[string]interface{}, the provided arguments should be key, value, key, value...")
	_ = Helpers().Register("dict", mapFn, "Creates a new map[string]interface{}, the provided arguments should be key, value, key, value...")
}

func defaultFn(val ...interface{}) interface{} {
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

func mapFn(args ...interface{}) map[string]interface{} {
	if len(args)%2 != 0 {
		log.Panicf("Error in 'map' directive. The number of keys do not match the number of values")
	}

	ret := map[string]interface{}{}

	for i := 0; i < len(args); i = i + 2 {
		ret[fmt.Sprintf("%v", args[i])] = args[i+1]
	}

	return ret
}
