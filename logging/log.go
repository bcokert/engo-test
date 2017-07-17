package logging

import (
	"fmt"
	"log"
	"os"
)

type LOGLEVEL int

const (
	DEBUG LOGLEVEL = iota
	INFO  LOGLEVEL = iota
)

const (
	colorinfo  = "\x1b[37m"
	colorerror = "\x1b[31m"
	colordebug = "\x1b[39m"
	endcolor   = "\x1b[0m"
)

// F is a shorthand for the fields struct accepted by log methods
type F map[string]interface{}

func colorize(color, text string) string {
	return color + text + endcolor
}

// A Logger is a structered logger, printing a message and zero or more named fields
type Logger interface {
	Error(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

type DefaultLogger struct {
	level  LOGLEVEL
	logger *log.Logger
}

func NewDefaultLogger(level LOGLEVEL) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *DefaultLogger) log(msg string, fields map[string]interface{}, color string, level string) {
	logmsg := fmt.Sprintf("%s %-50s\t", colorize(color, level), msg)
	for k, v := range fields {
		logmsg += fmt.Sprintf("%s=%+v\t", colorize(color, k), v)
	}
	l.logger.Println(logmsg)
}

func (l *DefaultLogger) Error(msg string, fields map[string]interface{}) {
	l.log(msg, fields, colorerror, "ERR")
}

func (l *DefaultLogger) Info(msg string, fields map[string]interface{}) {
	l.log(msg, fields, colorinfo, "INF")
}

func (l *DefaultLogger) Debug(msg string, fields map[string]interface{}) {
	if l.level == DEBUG {
		l.log(msg, fields, colordebug, "DBG")
	}
}
