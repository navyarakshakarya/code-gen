package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger provides structured logging with different levels
type Logger struct {
	verbose bool
}

// New creates a new logger instance
func New(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

// Info logs informational messages (only in verbose mode)
func (l *Logger) Info(format string, args ...interface{}) {
	if l.verbose {
		l.log("INFO", format, args...)
	}
}

// Success logs success messages
func (l *Logger) Success(format string, args ...interface{}) {
	l.log("✓", format, args...)
}

// Warning logs warning messages
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log("⚠", format, args...)
}

// Error logs error messages
func (l *Logger) Error(format string, args ...interface{}) {
	l.log("✗", format, args...)
}

// Fatal logs error and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log("✗", format, args...)
	os.Exit(1)
}

func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)

	// Use log package for consistent output
	log.Printf("[%s] %s %s", timestamp, level, message)
}
