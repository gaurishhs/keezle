package logger

import (
	"log"
)

type Logger interface {
	Log(msg string, v ...any)
}

type LogFn func(msg string, v ...any)

func (d LogFn) Log(msg string, v ...any) {
	d(msg, v...)
}

var (
	StdoutLogger = LogFn(func(msg string, v ...any) { log.Printf(msg, v...) })
	NoOpLogger   = LogFn(func(msg string, v ...any) {})
)
