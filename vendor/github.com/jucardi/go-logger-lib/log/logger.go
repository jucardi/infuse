package log

import "io"

// LoggerBuilder defines a logger constructor. The factory contains multiple logger constructors where new loggers with specified names
// can be created. Also these loggers can have their own io.Writer.
type LoggerBuilder func(name string, writer ...io.Writer) ILogger

// ILogger defines the contract for a logger interface to be used by the mgo and mongo packages.
// This interface matches most commonly used loggers which should make it simple to assign any
// logger implementation being used. By default it uses the sirupsen/logrus standard logger
// implementation.
type ILogger interface {
	// Name returns the manager name
	Name() string

	// SetLevel sets the logging level
	SetLevel(level Level)
	// GetLevel gets the logging level
	GetLevel() Level

	// Debug logs a message at level Debug on the logger.
	Debug(args ...interface{})
	// Debugf logs a message at level Debug on the logger.
	Debugf(format string, args ...interface{})

	// Info logs a message at level Info on the logger.
	Info(args ...interface{})
	// Infof logs a message at level Info on the logger.
	Infof(format string, args ...interface{})

	// Warn logs a message at level Warn on the logger.
	Warn(args ...interface{})
	// Warnf logs a message at level Warn on the logger.
	Warnf(format string, args ...interface{})

	// Error logs a message at level Error on the logger.
	Error(args ...interface{})
	// Errorf logs a message at level Error on the logger.
	Errorf(format string, args ...interface{})

	// Fatal logs a message at level Fatal on the logger.
	Fatal(args ...interface{})
	// Fatalf logs a message at level Fatal on the logger.
	Fatalf(format string, args ...interface{})

	// Panic logs a message at level Panic on the logger.
	Panic(args ...interface{})
	// Panicf logs a message at level Panic on the logger.
	Panicf(format string, args ...interface{})

	// SetFormatter sets a custom formatter to display the logs
	SetFormatter(formatter IFormatter)
}
