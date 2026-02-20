package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	green     = "\033[32m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	reset     = "\033[0m"
	infoLevel = "info"
	warnLevel = "warn"
	errLevel  = "error"
)

type Logger struct {
	level    string
	infoLog  *log.Logger
	// ErrorLog is exported so it can be assigned to http.Server.ErrorLog.
	ErrorLog *log.Logger
	warnLog  *log.Logger
}

func NewLogger(level string) Logger {
	logger := Logger{
		level:    level,
		infoLog:  log.New(os.Stdout, green+"INFO:\t"+reset, log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, red+"ERROR:\t"+reset, log.Ldate|log.Ltime),
		warnLog:  log.New(os.Stdout, yellow+"WARN:\t"+reset, log.Ldate|log.Ltime),
	}
	if logger.level == "" {
		logger.level = errLevel
	}
	logger.validateLogLevel()
	return logger
}

func (l *Logger) validateLogLevel() {
	if l.level == infoLevel || l.level == warnLevel || l.level == errLevel {
		return
	}
	l.warnLog.Printf("%s is not a valid log level. Log level must be '%s', '%s', or '%s'. Falling back to 'error' level",
		l.level, infoLevel, warnLevel, errLevel)
	l.level = errLevel
}

func (l *Logger) Info(msg string, args ...any) {
	if l.level == infoLevel {
		if len(args) != 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		l.infoLog.Printf("[%s] %s\n", logCaller(), msg)
	}
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.level == warnLevel || l.level == infoLevel {
		if len(args) != 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		l.warnLog.Printf("[%s] %s\n", logCaller(), msg)
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.ErrorLog.Printf("[%s] %s\n", logCaller(), msg)
}

func (l *Logger) ErrorErr(err error) {
	l.ErrorLog.Println("[" + logCaller() + "]" + err.Error())
}

func logCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
