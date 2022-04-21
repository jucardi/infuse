package log

import (
	"bytes"
	"io"
	"os"
)

type ILoggerAsync interface {
	ILogger
	Flush(clear ...bool)
	Reset()
}

type loggerMemory struct {
	ILogger
	buffer *bytes.Buffer
	writer io.Writer
}

// NewMemory returns a logger implementation that keeps the logs in memory until flushed. In this case,
// when providing an io.Writer, that would be the writer where all the stored logs would be flushed into.
// Uses os.Stdout if no writer is provided
func NewMemory(name string, writer ...io.Writer) ILogger {
	b := &bytes.Buffer{}
	l := NewLogrus(name, b).(ILogger)

	var w io.Writer = os.Stdout
	if len(writer) > 0 {
		w = writer[0]
	}

	return &loggerMemory{
		ILogger: l,
		buffer:  b,
		writer:  w,
	}
}

func (l *loggerMemory) Flush(clear ...bool) {
	_, _ = l.writer.Write(l.buffer.Bytes())
	if len(clear) > 0 && clear[0] {
		l.Reset()
	}
}

func (l *loggerMemory) Reset() {
	l.buffer.Reset()
}
