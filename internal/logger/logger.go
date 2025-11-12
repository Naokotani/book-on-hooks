package logger

import (
	"log"
	"os"
)

const (
	green     = "\033[32m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	reset     = "\033[0m"
	infoLevel = "info"
	warnLevel = "warn"
	errLevel  = "err"
)

type Logger struct {
	level    string
	infoLog  *log.Logger
	ErrorLog *log.Logger
	warnLog  *log.Logger
}

func NewLogger(level string) Logger {
	logger := Logger{
		level:    level,
		infoLog:  log.New(os.Stdout, green+"INFO:\t"+reset, log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, red+"ERROR:\t"+reset, log.Ldate|log.Ltime|log.Lshortfile),
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
		if len(args) == 0 {
			l.infoLog.Println(msg)
		} else {
			l.infoLog.Printf(msg, args...)
		}
	}
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.level == warnLevel || l.level == infoLevel {
		if len(args) == 0 {
			l.warnLog.Println(msg)
		} else {
			l.warnLog.Printf(msg, args...)
		}
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if len(args) == 0 {
		l.ErrorLog.Print(msg)
	} else {
		l.ErrorLog.Print(msg)
	}
}
