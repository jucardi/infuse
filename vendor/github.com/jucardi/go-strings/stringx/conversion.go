package stringx

import (
	"regexp"
	"strings"
)

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

// CamelToSnake converts CamelCase to snake_case.
func CamelToSnake(s string) string {
	return camelToSymbolSeparated(s, "_")
}

// CamelToDash converts CamelCase to dash-separated-string
func CamelToDash(s string) string {
	return camelToSymbolSeparated(s, "-")
}

// CamelToSnake converts CamelCase to snake_case.
func PascalToSnake(s string) string {
	return camelToSymbolSeparated(s, "_")
}

func PascalToDash(s string) string {
	return camelToSymbolSeparated(s, "-")
}

// SnakeToCamel converts snake_case to CamelCase
func SnakeToCamel(s string) string {
	return symbolSeparatedToCamel(s, "_")
}

// DashToCamel converts a dash-separated-string to CamelCase
func DashToCamel(s string) string {
	return symbolSeparatedToCamel(s, "-")
}

// SnakeToCamel converts snake_case to CamelCase
func SnakeToPascal(s string) string {
	return symbolSeparatedToPascal(s, "_")
}

// DashToCamel converts a dash-separated-string to CamelCase
func DashToPascal(s string) string {
	return symbolSeparatedToPascal(s, "-")
}

// ToTitle converts the provided space separated string into a title case string, ensuring the string is all lowercase but capitalizing the first letter of every word
func ToTitle(s string) string {
	words := New(s).ToLower().Split(" ")
	for i, v := range words {
		if len(v) == 0 {
			continue
		}
		words[i] = strings.ToUpper(v[:1]) + v[1:]
	}
	return strings.Join(words, " ")
}

func symbolSeparatedToPascal(s string, separator string) string {
	var ret string

	for _, v := range strings.Split(strings.ToLower(s), separator) {
		ret += strings.Title(v)
	}

	return ret
}

func symbolSeparatedToCamel(s string, separator string) string {
	str := symbolSeparatedToPascal(s, separator)
	return strings.ToLower(str[:1]) + str[1:]
}

func camelToSymbolSeparated(s string, separator string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, separator))
}
