package fmtc

import (
	"bytes"
	"fmt"
	"io"
)

type IFmt interface {
	Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)
	Sprintf(format string, a ...interface{}) string
	Errorf(format string, a ...interface{}) error
	Fprint(w io.Writer, a ...interface{}) (n int, err error)
	Print(a ...interface{}) (n int, err error)
	Sprint(a ...interface{}) string
	Fprintln(w io.Writer, a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Sprintln(a ...interface{}) string
}

type fmtFlow struct {
	colors []Color
}

// WithColors is an easy way to use the same function signature presented in the 'fmt' to print color messages.
//
//    Eg:   fmtc.WithColors(fmtc.White, fmtc.BgRed).Println(" danger! ")
//
func WithColors(colors ...Color) IFmt {
	return &fmtFlow{
		colors: colors,
	}
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, wrapColors(f.colors, format), a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(wrapColors(f.colors, format), a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func (f *fmtFlow) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(wrapColors(f.colors, format), a...)
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func (f *fmtFlow) Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(wrapColors(f.colors, format), a...)
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	doColors(f.colors, w)
	a = append(a, clear)
	return fmt.Fprint(w, a...)
}

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Print(a ...interface{}) (n int, err error) {
	doColors(f.colors)
	a = append(a, clear)
	return fmt.Print(a...)
}

// Sprint formats using the default formats for its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func (f *fmtFlow) Sprint(a ...interface{}) string {
	b := &bytes.Buffer{}
	doColors(f.colors, b)
	a = append(a, clear)
	return b.String() + fmt.Sprint(a...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	doColors(f.colors, w)
	a = append(a, clear)
	return fmt.Fprintln(w, a...)
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f *fmtFlow) Println(a ...interface{}) (n int, err error) {
	doColors(f.colors)
	a = append(a, clear)
	return fmt.Println(a...)
}

// Sprintln formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands and a newline is appended.
func (f *fmtFlow) Sprintln(a ...interface{}) string {
	b := &bytes.Buffer{}
	doColors(f.colors, b)
	a = append(a, clear)
	return b.String() + fmt.Sprintln(a...)
}
