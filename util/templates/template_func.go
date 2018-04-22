package templates

import (
	"fmt"
	"reflect"
	"text/template"

	"github.com/jucardi/gotmpl/util/log"
)

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"default":  defaultFn,
		"template": defaultFn,
		"map":      mapFn,
		"dict":     mapFn,
	}
}

func defaultFn(val interface{}) interface{} {
	log.Debug("--- templates.defaultFn")
	log.Debugf("------ return -> %v", val)
	return val
}

func templateFn(arg0 reflect.Value, args ...reflect.Value) reflect.Value {
	log.Debug("--- templates.templateFn")
	//log.Debug(arg0.Interface())
	return arg0
}

func mapFn(args ...interface{}) map[string]interface{} {
	log.Debug("--- templates.mapFn")

	if len(args)%2 != 0 {
		log.Panicf("Error in 'map' directive. The number of keys do not match the number of values")
	}

	ret := map[string]interface{}{}

	for i := 0; i < len(args); i = i + 2 {
		ret[fmt.Sprintf("%v", args[i])] = args[i+1]
	}

	log.Debug(ret)
	return ret
}
