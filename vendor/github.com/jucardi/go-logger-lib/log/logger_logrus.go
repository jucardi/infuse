package log

import (
	"bytes"
	"io"

	"github.com/sirupsen/logrus"
)

// LoggerLogrus indicates the name of the predefined logrus ILogger implementation
const LoggerLogrus = "logrus"

type logrusImpl struct {
	name string
	*logrus.Logger
}

func (l *logrusImpl) Name() string {
	return l.name
}

func (l *logrusImpl) SetLevel(level Level) {
	l.Level = logrus.Level(level)
}

func (l *logrusImpl) GetLevel() Level {
	return Level(l.Level)
}

func (l *logrusImpl) SetFormatter(formatter IFormatter) {
	l.Logger.SetFormatter(&logrusFormatter{
		l: l,
		f: formatter,
	})
}

// NewLogrus creates a new instance of the logrus implementation of ILogger
func NewLogrus(name string, writer ...io.Writer) ILogger {
	ret := &logrusImpl{
		name:   name,
		Logger: logrus.New(),
	}
	if len(writer) > 0 && writer[0] != nil {
		ret.Out = writer[0]
	}
	ret.SetFormatter(NewTerminalFormatter())
	return ret
}

type logrusFormatter struct {
	f IFormatter
	l ILogger
}

func (f *logrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := entry.Buffer
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}
	if err := f.f.Format(buffer, &Entry{
		LoggerName: f.l.Name(),
		Data:       entry.Data,
		Timestamp:  entry.Time,
		Level:      Level(entry.Level),
		Message:    entry.Message,
	}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
