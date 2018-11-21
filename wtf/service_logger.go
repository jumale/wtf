package wtf

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"path"
	"time"
)

const (
	LogLevelDebug LogLevel = 1 + iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

const levelMaxLength = 5

var logLevelMap = map[string]LogLevel{
	"debug":   LogLevelDebug,
	"info":    LogLevelInfo,
	"warn":    LogLevelWarn,
	"warning": LogLevelWarn,
	"error":   LogLevelError,
}

type LogLevel int

type Logger interface {
	Debug(msg string)
	Debugf(msg string, v ...interface{})
	Info(msg string)
	Infof(msg string, v ...interface{})
	Warn(msg string)
	Warnf(msg string, v ...interface{})
	Error(msg string)
	Errorf(msg string, v ...interface{})
	Close() error
}

func NewFileLogger(configDir string, cfg AppConfig) (*FileLogger, error) {
	var err error
	logger := FileLogger{
		filePath:   path.Join(configDir, cfg.Log.File),
		dateFormat: cfg.Log.DateFormat,
		level:      cfg.Log.Level,
	}

	logger.file, err = os.OpenFile(logger.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &logger, nil
}

type FileLogger struct {
	filePath   string
	file       *os.File
	dateFormat string
	level      LogLevel
}

func (l FileLogger) Debug(msg string) {
	if l.level <= LogLevelDebug {
		l.log(msg, "DEBUG")
	}
}

func (l FileLogger) Debugf(msg string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log(fmt.Sprintf(msg, v...), "DEBUG")
	}
}

func (l FileLogger) Info(msg string) {
	if l.level <= LogLevelInfo {
		l.log(msg, "INFO")
	}
}

func (l FileLogger) Infof(msg string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log(fmt.Sprintf(msg, v...), "INFO")
	}
}

func (l FileLogger) Warn(msg string) {
	if l.level <= LogLevelWarn {
		l.log(msg, "WARG")
	}
}

func (l FileLogger) Warnf(msg string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.log(fmt.Sprintf(msg, v...), "WARN")
	}
}

func (l FileLogger) Error(msg string) {
	if l.level <= LogLevelError {
		l.log(msg, "ERROR")
	}
}

func (l FileLogger) Errorf(msg string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.log(fmt.Sprintf(msg, v...), "ERROR")
	}
}

func (l FileLogger) Close() error {
	return l.file.Close()
}

func (l FileLogger) log(msg, level string) {
	pad := fmt.Sprintf("-%d", levelMaxLength+2)
	level = fmt.Sprintf("[%s]", level)
	msg = fmt.Sprintf("%"+pad+"s %s", level, msg)

	if l.dateFormat != "" {
		t := time.Now()
		msg = fmt.Sprintf("[%s] %s", t.Format(l.dateFormat), msg)
	}

	_, err := l.file.Write([]byte(msg + "\n"))
	if err != nil {
		log.Panicf("counld not write log message to file '%s'", l.filePath)
	}
}

func (ll *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var label string
	err := unmarshal(&label)
	if err != nil {
		return errors.WithStack(err)
	}

	var ok bool
	if *ll, ok = logLevelMap[label]; !ok {
		return errors.Errorf("wrong config: expected log level one of %v, actual '%s'", allowedLogLevels(), label)
	}
	return nil
}

func allowedLogLevels() (levels []string) {
	for label := range logLevelMap {
		levels = append(levels, label)
	}
	return levels
}
