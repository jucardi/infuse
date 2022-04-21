package stringx

import (
	"strings"
	"unicode"
)

type String struct {
	current string
}

func New(str string) String {
	return String{current: str}
}

func Join(a []string, sep string) String {
	return New(strings.Join(a, sep))
}

// Replace returns a copy of the string s with the first n
// non-overlapping instances of old replaced by new.
func (s String) Replace(old, new string, n int) String {
	return New(strings.Replace(s.current, old, new, n))
}

// Repeat returns a new string consisting of count copies of the string s.
//
// It panics if count is negative or if
// the result of (len(s) * count) overflows.
func (s String) Repeat(count int) String {
	return New(strings.Repeat(s.current, count))
}

// ToUpper returns a copy of the string s with all Unicode letters mapped to their upper case.
func (s String) ToUpper() String {
	return New(strings.ToUpper(s.current))
}

// ToLower returns a copy of the string s with all Unicode letters mapped to their lower case.
func (s String) ToLower() String {
	return New(strings.ToLower(s.current))
}

// ToTitle returns a copy of the string s with all Unicode letters mapped to their title case.
func (s String) ToTitle() String {
	return New(strings.ToTitle(s.current))
}

// ToUpperSpecial returns a copy of the string s with all Unicode letters mapped to their
// upper case, giving priority to the special casing rules.
func (s String) ToUpperSpecial(c unicode.SpecialCase) String {
	return New(strings.ToUpperSpecial(c, s.current))
}

// ToLowerSpecial returns a copy of the string s with all Unicode letters mapped to their
// lower case, giving priority to the special casing rules.
func (s String) ToLowerSpecial(c unicode.SpecialCase) String {
	return New(strings.ToLowerSpecial(c, s.current))
}

// ToTitleSpecial returns a copy of the string s with all Unicode letters mapped to their
// title case, giving priority to the special casing rules.
func (s String) ToTitleSpecial(c unicode.SpecialCase) String {
	return New(strings.ToTitleSpecial(c, s.current))
}

// Trim returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutset removed.
func (s String) Trim(cutset string) String {
	return New(strings.Trim(s.current, cutset))
}

// TrimLeft returns a slice of the string s with all leading
// Unicode code points contained in cutset removed.
func (s String) TrimLeft(cutset string) String {
	return New(strings.TrimLeft(s.current, cutset))
}

// TrimRight returns a slice of the string s, with all trailing
// Unicode code points contained in cutset removed.
func (s String) TrimRight(cutset string) String {
	return New(strings.TrimRight(s.current, cutset))
}

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func (s String) TrimSpace() String {
	return New(strings.TrimSpace(s.current))
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func (s String) TrimPrefix(prefix string) String {
	return New(strings.TrimPrefix(s.current, prefix))
}

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func (s String) TrimSuffix(suffix string) String {
	return New(strings.TrimSuffix(s.current, suffix))
}

// Title returns a copy of the string s with all Unicode letters that begin words
// mapped to their title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
func (s String) Title() String {
	return New(strings.Title(s.current))
}

// Map returns a copy of the string s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the string with no replacement.
func (s String) Map(mapping func(rune) rune) String {
	return New(strings.Map(mapping, s.current))
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func (s String) Index(sep string) int {
	return strings.Index(s.current, sep)
}

// Count counts the number of non-overlapping instances of sep in s.
// If sep is an empty string, Count returns 1 + the number of Unicode code points in s.
func (s String) Count(sep string) int {
	return strings.Count(s.current, sep)
}

// Contains reports whether substr is within s.
func (s String) Contains(substr string) bool {
	return strings.Contains(s.current, substr)
}

// ContainsAny reports whether any Unicode code points in chars are within s.
func (s String) ContainsAny(chars string) bool {
	return strings.ContainsAny(s.current, chars)
}

// ContainsRune reports whether the Unicode code point r is within s.
func (s String) ContainsRune(r rune) bool {
	return strings.ContainsRune(s.current, r)
}

// HasPrefix tests whether the string s begins with prefix.
func (s String) HasPrefix(prefix string) bool {
	return strings.HasPrefix(s.current, prefix)
}

// HasSuffix tests whether the string s ends with suffix.
func (s String) HasSuffix(suffix string) bool {
	return strings.HasSuffix(s.current, suffix)
}

// SplitN slices s into substrings separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, SplitN splits after each UTF-8 sequence.
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings
func (s String) SplitN(sep string, n int) []string {
	return strings.SplitN(s.current, sep, n)
}

// SplitAfterN slices s into substrings after each instance of sep and
// returns a slice of those substrings.
// If sep is empty, SplitAfterN splits after each UTF-8 sequence.
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings
func (s String) SplitAfterN(sep string, n int) []string {
	return strings.SplitAfterN(s.current, sep, n)
}

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
// If sep is empty, Split splits after each UTF-8 sequence.
// It is equivalent to SplitN with a count of -1.
func (s String) Split(sep string) []string {
	return strings.Split(s.current, sep)
}

// SplitAfter slices s into all substrings after each instance of sep and
// returns a slice of those substrings.
// If sep is empty, SplitAfter splits after each UTF-8 sequence.
// It is equivalent to SplitAfterN with a count of -1.
func (s String) SplitAfter(sep string) []string {
	return strings.SplitAfter(s.current, sep)
}

// FieldsFunc splits the string s at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of s. If all code points in s satisfy f(c) or the
// string is empty, an empty slice is returned.
// FieldsFunc makes no guarantees about the order in which it calls f(c).
// If f does not return consistent results for a given c, FieldsFunc may crash.
func (s String) FieldsFunc(f func(rune) bool) []string {
	return strings.FieldsFunc(s.current, f)
}

// Fields splits the string s around each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, returning an array of substrings of s or an
// empty list if s contains only white space.
func (s String) Fields() []string {
	return strings.Fields(s.current)
}

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding.
func (s String) EqualFold(t string) bool {
	return strings.EqualFold(s.current, t)
}

// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
func (s String) IndexByte(c byte) int {
	return strings.IndexByte(s.current, c)
}

// IndexFunc returns the index into s of the first Unicode
// code point satisfying f(c), or -1 if none do.
func (s String) IndexFunc(f func(rune) bool) int {
	return strings.IndexFunc(s.current, f)
}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the first instance of any
// invalid UTF-8 byte sequence.
func (s String) IndexRune(r rune) int {
	return strings.IndexRune(s.current, r)
}

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
func (s String) IndexAny(chars string) int {
	return strings.IndexAny(s.current, chars)
}

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.
func (s String) LastIndex(sep string) int {
	return strings.LastIndex(s.current, sep)
}

// LastIndexAny returns the index of the last instance of any Unicode code
// point from chars in s, or -1 if no Unicode code point from chars is
// present in s.
func (s String) LastIndexAny(chars string) int {
	return strings.LastIndexAny(s.current, chars)
}

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func (s String) LastIndexByte(c byte) int {
	return strings.LastIndexByte(s.current, c)
}

// LastIndexFunc returns the index into s of the last
// Unicode code point satisfying f(c), or -1 if none do.
func (s String) LastIndexFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(s.current, f)
}

// NewReader returns a new Reader reading from s.
// It is similar to bytes.NewBufferString but more efficient and read-only.
func (s String) NewReader() *strings.Reader {
	return strings.NewReader(s.current)
}

// S returns the current state of the string
func (s String) S() string {
	return s.current
}
