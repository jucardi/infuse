package ioutils

import "github.com/jucardi/go-strings/stringx"

type StringWriter struct {
	buffer *stringx.StringBuilder
}

func (s *StringWriter) Write(p []byte) (n int, err error) {
	if s.buffer == nil {
		s.buffer = stringx.Builder()
	}

	s.buffer.Append(string(p))
	return len(p), nil
}

func (s *StringWriter) ToString() string {
	if s.buffer == nil {
		return ""
	}

	return s.buffer.Build()
}

func NewStringWriter() *StringWriter {
	return &StringWriter{}
}
