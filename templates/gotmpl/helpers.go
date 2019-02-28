package gotmpl

import (
	"fmt"
	"github.com/jucardi/infuse/templates/helpers"
	"github.com/jucardi/infuse/util/log"
	"os"
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
	Helpers().Register("default", defaultFn)
	Helpers().Register("map", mapFn)
	Helpers().Register("dict", mapFn)
	Helpers().Register("env", envFn)
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

func envFn(name string) string {
	return os.Getenv(name)
}
