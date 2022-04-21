package log

import "github.com/jucardi/go-logger-lib/log"

// For more information about how the loggin works in the github.com/jucardi/go-logger-lib/log
// package, please refer to https://github.com/jucardi/go-logger-lib/blob/master/README.md

// LoggerOsx defines the name for the logger used for the mgo package
const LoggerOsx = "osx"

// Get returns the current logger instance
func Get() log.ILogger {
	return log.Get(LoggerOsx)
}

// Set assigns a new ILogger instance to be used as the logger for the mgo package
func Set(logger log.ILogger) {
	log.Register(LoggerOsx, logger)
}

// Disable disables logging for the mgo package by assigning a the nil implementation of ILogger
func Disable() {
	Set(log.NewNil())
}
