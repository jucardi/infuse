package stringx

// GetOrDefault is a helper for variadic arguments when used as single optional input value in a variadic function.
func GetOrDefault(def string, args ...string) string {
	if len(args) > 0 {
		return args[0]
	}
	return def
}
