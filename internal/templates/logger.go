package templates

// LoggerTemplate remains unchanged
const LoggerTemplate = `package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger interface
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
}

var log *logrus.Logger

// Init initializes the logger
func Init(level string) {
	log = logrus.New()
	
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	
	log.SetLevel(logLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

// GetLogger returns the logger instance
func GetLogger() Logger {
	return &logrusLogger{}
}

type logrusLogger struct{}

// Info logs info level message
func (l *logrusLogger) Info(args ...interface{}) {
	log.Info(args...)
}

// Error logs error level message
func (l *logrusLogger) Error(args ...interface{}) {
	log.Error(args...)
}

// Debug logs debug level message
func (l *logrusLogger) Debug(args ...interface{}) {
	log.Debug(args...)
}

// Warn logs warn level message
func (l *logrusLogger) Warn(args ...interface{}) {
	log.Warn(args...)
}

// Package level functions for backward compatibility
func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}
`
