package stringx

import (
	"bytes"
	"fmt"
	"strconv"
)

// StringBuilder encapsulates the buffer to use for the builder
type StringBuilder struct {
	buffer bytes.Buffer
}

// Builder Creates a new StringBuilder
func Builder() *StringBuilder {
	return &StringBuilder{}
}

// Append Appends the given string(s) to the builder
func (s *StringBuilder) Append(args ...string) *StringBuilder {
	for _, v := range args {
		s.buffer.WriteString(v)
	}
	return s
}

// Appendf processes a string format based on the format and arguments passed and appends the result to the builder
func (s *StringBuilder) Appendf(format string, args ...interface{}) *StringBuilder {
	return s.Append(fmt.Sprintf(format, args...))
}

// AppendObj attempts to obtain a string representation of the interface{} (s) and appends it/them to the builder
func (s *StringBuilder) AppendObj(objs ...interface{}) *StringBuilder {
	for _, v := range objs {
		s.Append(fmt.Sprintf("%+v", v))
	}
	return s
}

// AppendInt Appends the given number(s) to the builder
func (s *StringBuilder) AppendInt(i ...int) *StringBuilder {
	for _, v := range i {
		s.Append(strconv.Itoa(v))
	}
	return s
}

// AppendLine Appends the given string(s) to the builder.
func (s *StringBuilder) AppendLine(lines ...string) *StringBuilder {
	for _, v := range lines {
		s.Append(v).Br()
	}
	return s
}

// AppendLinef processes a string format based on the format and arguments passed, appends the result to the builder and a new line at the end
func (s *StringBuilder) AppendLinef(format string, args ...interface{}) *StringBuilder {
	return s.AppendLine(fmt.Sprintf(format, args...))
}

// AppendRune Appends a single character to the builder
func (s *StringBuilder) AppendRune(char rune) *StringBuilder {
	s.buffer.WriteRune(char)
	return s
}

// Br Breaks to the next line
func (s *StringBuilder) Br() *StringBuilder {
	return s.Append(LineBreak)
}

// IsEmpty Indicates whether the builder is empty
func (s *StringBuilder) IsEmpty() bool {
	return s.buffer.Len() == 0
}

// Build Builds the string
func (s *StringBuilder) Build() string {
	return s.buffer.String()
}
