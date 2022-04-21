package fmtc

import (
	"bytes"
	"fmt"
	"io"
)

const (
	esc   = "\033["
	clear = esc + "0m"
)

type IColorFlow interface {
	String() string
	Print(str interface{}, colors ...Color) IColorFlow
	PrintLn(str interface{}, colors ...Color) IColorFlow
	Printf(format string, args Args, colors ...Color) IColorFlow
}

type colorFlow struct {
	customWriter bool
	writer       io.Writer
}

// Args is an alias of []interface{}
type Args []interface{}

// New creates a new IColorFlow
//
// - writer: (Optional) An io.Writer where the messages will be written. If using a writer, `String()` is not supported.
//
func New(writer ...io.Writer) IColorFlow {
	ret := &colorFlow{}

	if len(writer) > 0 && writer[0] != nil {
		ret.writer = writer[0]
		ret.customWriter = true
	} else {
		ret.writer = &bytes.Buffer{}
	}

	return ret
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, str interface{}, colors ...Color) {
	doColors(colors, w)
	fmt.Fprint(w, str, clear)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Fprintln(w io.Writer, str interface{}, colors ...Color) {
	doColors(colors, w)
	fmt.Fprintln(w, str, clear)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, args Args, colors ...Color) {
	doColors(colors, w)
	fmt.Fprintf(w, format, args...)
	fmt.Fprint(w, clear)
}

// String returns the resulting string of the color flow.
func (f *colorFlow) String() string {
	if f.customWriter {
		return ""
	}
	buf := f.writer.(*bytes.Buffer)
	return buf.String()
}

// Print prints the message using the provided colors into the writer used by the color flow
func (f *colorFlow) Print(str interface{}, colors ...Color) IColorFlow {
	Fprint(f.writer, str, colors...)
	return f
}

// Print prints the message using the provided colors into the writer used by the color flow, adds a line break at the end
func (f *colorFlow) PrintLn(str interface{}, colors ...Color) IColorFlow {
	Fprintln(f.writer, str, colors...)
	return f
}

// Printf prints the message by doing a string format with the provided `fmtc.Args`. Prints the result of the format using the
// provided colors into the writer used by the color flow.
func (f *colorFlow) Printf(format string, args Args, colors ...Color) IColorFlow {
	Fprintf(f.writer, format, args, colors...)
	return f
}
