package log

import (
	"io"
	"time"
)

const (
	FieldLevel      = "level"
	FieldTimestamp  = "timestamp"
	FieldMessage    = "message"
	FieldLoggerName = "loggerName"
)

// Entry represents a log entry.
type Entry struct {
	// LoggerName indicates to what logger the log entry belongs to
	LoggerName string

	// Contains all the fields set by the user. TODO
	Data map[string]interface{}

	// Time at which the log entry was created
	Timestamp time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	metadata map[string]interface{}
	writer   io.Writer
}

func (e *Entry) AddMetadata(key string, val interface{}) {
	if e.metadata == nil {
		e.metadata = map[string]interface{}{}
	}
	e.metadata[key] = val
}

func (e *Entry) getField(name string) interface{} {
	switch name {
	case FieldMessage:
		return e.Message
	case FieldLevel:
		return e.Level
	case FieldTimestamp:
		return e.Timestamp
	case FieldLoggerName:
		return e.LoggerName
	}
	if e.Data == nil {
		return nil
	}
	if ret, ok := e.Data[name]; ok {
		return ret
	}
	return nil
}
