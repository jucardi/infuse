package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/jucardi/infuse/util/log"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

/** In this file are defined the generic helpers that may work for different template types. */

// RegisterCommon registers the generic helpers designed to work for different template types.
func RegisterCommon(manager IHelpersManager) {
	_ = manager.Register("string", stringFn, "Prints a string representation of the provided object")
	_ = manager.Register("format", fmt.Sprintf, "Formats according to a format specifier and returns the resulting string.")
	_ = manager.Register("uppercase", strings.ToUpper, "Returns a copy of the string s with all Unicode letters mapped to their upper case.")
	_ = manager.Register("lowercase", strings.ToLower, "Returns a copy of the string s with all Unicode letters mapped to their lower case.")
	_ = manager.Register("title", strings.ToTitle, "Returns a copy of the string s with all Unicode letters mapped to their title case.")
	_ = manager.Register("stringsReplace", strings.Replace, "Returns a copy of the string s with the first n non-overlapping instances of old replaced by new.")
	_ = manager.Register("stringsJoin", strings.Join, "Concatenates the elements of a to create a single string. The separator string sep is placed between elements in the resulting string.")
	_ = manager.Register("stringsSplit", strings.Split, "Slices s into all substrings separated by sep and returns a slice of the substrings between those separators.")
	_ = manager.Register("stringsTrim", strings.Trim, "Returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimLeft", strings.TrimLeft, "Returns a slice of the string s with all leading Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimRight", strings.TrimRight, "Returns a slice of the string s, with all trailing Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimSpace", strings.TrimSpace, "Returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode.")
	_ = manager.Register("startsWith", strings.HasPrefix, "Returns a boolean indicating whether the string s begins with prefix.")
	_ = manager.Register("endsWith", strings.HasSuffix, "Returns a boolean indicating whether the string s ends with suffix.")
	_ = manager.Register("br", bracketsFn, "Wraps the contents into double brackets {{ }}")
	_ = manager.Register("yaml", toYMLString, "Marshals the provided object as YAML")
	_ = manager.Register("json", toJSONString, "Marshals the provided object as JSON")
	_ = manager.Register("rem", comment, "Helper to add comments")
	_ = manager.Register("env", os.Getenv, "Returns the value set in the provided environment variable")
	_ = manager.Register("stringArray", stringArray, "Creates an array of strings with the provided string args")
}

/** String helpers */

func stringFn(arg interface{}) string {
	return fmt.Sprintf("%+v", arg)
}

func bracketsFn(arg interface{}) string {
	return fmt.Sprintf("{{%+v}}", arg)
}

func toYMLString(arg interface{}) string {
	result, err := yaml.Marshal(arg)
	if err != nil {
		log.Panicf("error occurred while converting object to YAML. %+v", err)
	}
	return string(result)
}

func toJSONString(arg interface{}) string {
	result, err := json.Marshal(arg)
	if err != nil {
		log.Panicf("error occurred while converting object to JSON. %+v", err)
	}
	return string(result)
}

func comment(_ ...interface{}) string {
	return ""
}

func stringArray(args ...string) []string {
	return args
}