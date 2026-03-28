// Package log is a basic logger designed to be simple and readable
package log

import "fmt"

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var labels = map[int]string{INFO: "INFO", DEBUG: "DEBUG", WARN: "WARN", ERROR: "ERROR"}

type Logger struct {
	level int
	label string
}

func DefaultLogger() *Logger {
	return &Logger{
		level: DEBUG,
		label: "LOG",
	}
}

func NewLogger(level int, label string) *Logger {
	return &Logger{
		level: level,
		label: label,
	}
}

func (l *Logger) Log(level int, format string, values ...any) {
	if level == l.level || level < l.level {
		fmt.Printf("[%v - %v]: %v", l.label, labels[level], fmt.Sprintf(format, values...))
	}
}

func (l *Logger) Info(format string, values ...any) {
	l.Log(INFO, format, values...)
}

func (l *Logger) Debug(format string, values ...any) {
	l.Log(DEBUG, format, values...)
}

func (l *Logger) Warn(format string, values ...any) {
	l.Log(WARN, format, values...)
}

func (l *Logger) Error(format string, values ...any) {
	l.Log(ERROR, format, values...)
}
