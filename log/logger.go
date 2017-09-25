package log

import (
	"fmt"
	"os"
)

func New() Logger {
	return &logger{}
}

// Logger - simple interface for write to stdout and stderr
type Logger interface {
	Printf(format string, args ...interface{})
	Error(args ...interface{})
}

type NoopLogger struct{}

func (NoopLogger) Printf(string, ...interface{}) {}

func (NoopLogger) Error(...interface{}) {}

type logger struct {
}

func (l *logger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

func (l *logger) Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
