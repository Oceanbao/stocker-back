package infra

import (
	"log/slog"
)

type LoggerSlog struct {
	logger *slog.Logger
}

func NewLoggerSlog(logger *slog.Logger) *LoggerSlog {
	return &LoggerSlog{
		logger: logger,
	}
}

func (l *LoggerSlog) Debugf(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *LoggerSlog) Infof(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *LoggerSlog) Warnf(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *LoggerSlog) Errorf(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}
