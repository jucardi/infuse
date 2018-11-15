package log

import (
	"github.com/jucardi/infuse/config"
	"github.com/sirupsen/logrus"
	"os"
)

// Debug logs a message at the debug level in the standard logrus logger
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf logs a message at the debug level in the standard logrus logger
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info logs a message at the info level in the standard logrus logger
func Info(args ...interface{}) {
	logrus.Debug(args...)
}

// Infof logs a message at the info level in the standard logrus logger
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Error logs a message at the error level in the standard logrus logger
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf logs a message at the error level in the standard logrus logger
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Panic logs a message at the error level in the standard logrus logger and exits the application. Used logrus.Error to avoid printing unnecessary stack trace.
func Panic(args ...interface{}) {
	logrus.Error(args...)
	if config.Get().Verbose {
		// logrus panic adds multiple unnecesary internal logrus traces to the stack.
		panic("")
	}
	os.Exit(-1)
}

// Panicf logs a message at the error level in the standard logrus logger and exits the application. Used logrus.Errorf to avoid printing unnecessary stack trace.
func Panicf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
	if config.Get().Verbose {
		// logrus panic adds multiple unnecesary internal logrus traces to the stack.
		panic("")
	}

	os.Exit(-1)
}
