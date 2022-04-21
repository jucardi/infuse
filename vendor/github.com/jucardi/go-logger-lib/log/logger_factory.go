package log

import (
	"io"
)

var (
	loggers = map[string]ILogger{}
)

// Register registers an instance of ILogger to be returned as the singleton
// instance by the given name.
//
//   {name}   - The logger name.
//   {logger} - The logger instance.
//
func Register(name string, logger ILogger) ILogger {
	loggers[name] = logger
	return logger
}

// Get returns an instance of the requested logger by its name. Returns the Nil Logger implementation
// if a logger by the given name is not found.
//
//   {name} - The name of the logger instance to be retrieved.
//
func Get(name string) ILogger {
	if v, ok := loggers[name]; ok {
		return v
	}

	return Register(name, defaultBuilder(name))
}

// New creates a new logger instance using the default builder assigned.
//
//   {name}   - The name of the logger to create.
//   {writer} - (Optional) The io.Writer the logger instance should use. If not provided,
//              it is set to the default writer by the implementation, typically Stdout or Stderr
//
func New(name string, writer ...io.Writer) ILogger {
	return Register(name, defaultBuilder(name, writer...))
}

// List returns the list of loggers that have been registered.
func List() []string {
	var ret []string
	for k := range loggers {
		ret = append(ret, k)
	}
	return ret
}

// SetDefaultBuilder assigns the default builder to be used when creating new loggers.
func SetDefaultBuilder(ctor LoggerBuilder) {
	defaultBuilder = ctor
}

// Contains indicates if a logger by the given name exists.
func Contains(name string) bool {
	for k := range loggers {
		if name == k {
			return true
		}
	}

	return false
}
