package logger

import (
	"log"
)

// Logger is an interface for logging messages.
type Logger interface {
	Log(msg string, v ...any)
}

// LogFn is a function type that implements the Logger interface.
type LogFn func(msg string, v ...any)

// Log implements the Logger interface for LogFn.
func (d LogFn) Log(msg string, v ...any) {
	d(msg, v...)
}

// StdoutLogger is a logger that writes to standard output.
// NoLogger is a logger that does nothing.
var (
	StdoutLogger = LogFn(func(msg string, v ...any) { log.Printf(msg, v...) })
	NoLogger     = LogFn(func(msg string, v ...any) {})
)
