// Package logger provides a small set of leveled loggers (info, warning and
// error) that are shared across the application packages (main, handlers, db,
// forms and models).
//
// Every logger writes the log level, the date, the time and the affected file
// with its line number, for example:
//
//	INFO     2009/01/23 01:23:23 main.go:42: Server running...
//	WARNING  2009/01/23 01:23:23 foods.go:65: no rows affected
//	ERROR    2009/01/23 01:23:23 users.go:30: duplicate email
package logger

import (
	"io"
	"log"
	"os"
)

// Logger groups the three leveled loggers used throughout the application.
type Logger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// standard flags shared by every logger: date, time and the short file name
// with its line number.
const flags = log.Ldate | log.Ltime | log.Lshortfile

// New returns a Logger writing INFO and WARNING messages to infoOut and ERROR
// messages to errorOut. Passing nil for either writer falls back to os.Stdout
// (for info/warning) and os.Stderr (for errors).
func New(infoOut, errorOut io.Writer) *Logger {
	if infoOut == nil {
		infoOut = os.Stdout
	}
	if errorOut == nil {
		errorOut = os.Stderr
	}

	return &Logger{
		Info:    log.New(infoOut, "INFO\t", flags),
		Warning: log.New(infoOut, "WARNING\t", flags),
		Error:   log.New(errorOut, "ERROR\t", flags),
	}
}

// Default returns a Logger writing info/warning to stdout and errors to stderr.
func Default() *Logger {
	return New(os.Stdout, os.Stderr)
}

// Discard returns a Logger that throws away all output. It is handy for tests
// and as a nil-safe fallback when a component has not been given a logger.
func Discard() *Logger {
	return New(io.Discard, io.Discard)
}
