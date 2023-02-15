package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/infuse/util/log"
	"gopkg.in/yaml.v2"
)

/** In this file are defined the generic helpers that may work for different template types. */

// RegisterCommon registers the generic helpers designed to work for different template types.
func RegisterCommon(manager IHelpersManager) {
	_ = manager.Register("string", stringFn, "Prints a string representation of the provided object")
	_ = manager.Register("format", fmt.Sprintf, "Formats according to a format specifier and returns the resulting string.")
	_ = manager.Register("uppercase", strings.ToUpper, "Returns a copy of the string s with all Unicode letters mapped to their upper case.")
	_ = manager.Register("upper", strings.ToUpper, "Returns a copy of the string s with all Unicode letters mapped to their upper case.")
	_ = manager.Register("lowercase", strings.ToLower, "Returns a copy of the string s with all Unicode letters mapped to their lower case.")
	_ = manager.Register("lower", strings.ToLower, "Returns a copy of the string s with all Unicode letters mapped to their lower case.")
	_ = manager.Register("title", stringx.ToTitle, "Returns a copy of the string s with all Unicode letters mapped to their title case.")
	_ = manager.Register("stringsReplace", strings.Replace, "Returns a copy of the string s with the first n non-overlapping instances of old replaced by new.")
	_ = manager.Register("stringsJoin", stringsJoin, "Concatenates the elements of a to create a single string. The separator string sep is placed between elements in the resulting string.")
	_ = manager.Register("stringsSplit", strings.Split, "Slices s into all substrings separated by sep and returns a slice of the substrings between those separators.")
	_ = manager.Register("stringsTrim", strings.Trim, "Returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimLeft", strings.TrimLeft, "Returns a slice of the string s with all leading Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimRight", strings.TrimRight, "Returns a slice of the string s, with all trailing Unicode code points contained in cutset removed.")
	_ = manager.Register("stringsTrimSpace", strings.TrimSpace, "Returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode.")
	_ = manager.Register("stringsContains", strings.Contains, "Returns a boolean indicating whether the string s contains substr.")
	_ = manager.Register("stringsCompare", strings.Compare, "Returns an integer comparing two strings lexicographically.")
	_ = manager.Register("stringsSub", stringsSub, "Returns a substring of the specified string. E.g: {{stringsSub $sourceStr, startIndex, endIndex}}")
	_ = manager.Register("startsWith", strings.HasPrefix, "Returns a boolean indicating whether the string s begins with prefix.")
	_ = manager.Register("endsWith", strings.HasSuffix, "Returns a boolean indicating whether the string s ends with suffix.")
	_ = manager.Register("br", bracketsFn, "Wraps the contents into double brackets {{ }}")
	_ = manager.Register("yaml", toYMLString, "Marshals the provided object as YAML")
	_ = manager.Register("json", toJSONString, "Marshals the provided object as JSON")
	_ = manager.Register("rem", comment, "Helper to add comments")
	_ = manager.Register("env", os.Getenv, "Returns the value set in the provided environment variable")
	_ = manager.Register("stringArray", stringArray, "Creates an array of strings with the provided string args")
	_ = manager.Register("mathAdd", mathAdd, "Adds all the provided numbers together and returns the result")
	_ = manager.Register("mathMult", mathMult, "Multiplies all the provided numbers together and returns the result")
	_ = manager.Register("indent", indent, "Indents a given string using the provided indentation")
	_ = manager.Register("stringxCamelToDash", stringx.CamelToDash, "Converts a camelCase string to a dash-separated string")
	_ = manager.Register("stringxCamelToSnake", stringx.CamelToSnake, "Converts a camelCase string to a snake_case string")
	_ = manager.Register("stringxDashToCamel", stringx.DashToCamel, "Converts a dash-separated string to a camelCase string")
	_ = manager.Register("stringxSnakeToCamel", stringx.SnakeToCamel, "Converts a SnakeToCamel string to a camelCase string")
	_ = manager.Register("stringxPascalToDash", stringx.PascalToDash, "Converts a PascalCase string to a dash-separated string")
	_ = manager.Register("stringxPascalToSnake", stringx.PascalToSnake, "Converts a PascalCase string to a snake_case string")
	_ = manager.Register("stringxDashToPascal", stringx.DashToPascal, "Converts a dash-separated string to a PascalCase string")
	_ = manager.Register("stringxSnakeToPascal", stringx.SnakeToPascal, "Converts a SnakeToCamel string to a PascalCase string")
}

/** String helpers */

func stringsJoin(array interface{}, sep string) string {
	if kind := reflect.TypeOf(array).Kind(); kind != reflect.Array && kind != reflect.Slice {
		log.Panic("the first argument must be a valid array", kind.String())
	}
	val := reflect.ValueOf(array)
	if !val.IsValid() {
		return ""
	}
	arr := streams.From(array).Map(func(i interface{}) interface{} {
		if str, ok := i.(string); ok {
			return str
		}
		return fmt.Sprint(i)
	}).ToArray().([]string)
	return strings.Join(arr, sep)
}

func stringFn(arg interface{}) string {
	return fmt.Sprintf("%+v", arg)
}

func stringsSub(sourceStr string, start, end int) string {
	return sourceStr[start:end]
}

func bracketsFn(arg interface{}) string {
	return fmt.Sprintf("{{%+v}}", arg)
}

func toYMLString(arg interface{}, indent ...int) string {
	result, err := yaml.Marshal(arg)
	if err != nil {
		log.Panicf("error occurred while converting object to YAML. %+v", err)
	}
	ret := string(result)
	if len(indent) > 0 {
		ind := ""
		for i := 0; i < indent[0]; i++ {
			ind = ind + " "
		}
		split := strings.Split(ret, "\n")
		for i := 0; i < len(split); i++ {
			split[i] = ind + split[i]
		}
		ret = strings.Join(split, "\n")
	}
	return ret
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

func mathAdd(nums ...int) int {
	ret := 0
	for _, n := range nums {
		ret += n
	}
	return ret
}

func mathMult(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	ret := 1
	for _, n := range nums {
		ret *= n
	}
	return ret
}

func indent(str string, indent string) string {
	return strings.Replace(str, "\n", "\n"+indent, -1)
}
