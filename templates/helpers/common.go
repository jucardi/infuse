package helpers

import (
	"fmt"
	"strings"
)

/** In this file are defined the generic helpers that may work for different template types. */

// RegisterCommon registers the generic helpers designed to work for different template types.
func RegisterCommon(manager IHelpersManager) {
	manager.Register("string", stringFn)
	manager.Register("format", formatFn)
	manager.Register("uppercase", uppercaseFn)
	manager.Register("lowercase", lowercaseFn)
	manager.Register("br", bracketsFn)
}

/** String helpers */

func stringFn(arg interface{}) string {
	return fmt.Sprintf("%+v", arg)
}

func formatFn(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func lowercaseFn(arg interface{}) string {
	return strings.ToLower(stringFn(arg))
}

func uppercaseFn(arg string) string {
	return strings.ToUpper(stringFn(arg))
}

func bracketsFn(arg interface{}) string {
	return fmt.Sprintf("{{%+v}}", arg)
}
