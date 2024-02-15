package infra

import (
	"github.com/rs/zerolog"
)

type LoggerZero struct {
	logger *zerolog.Logger
}

func NewLoggerZero(logger *zerolog.Logger) *LoggerZero {
	return &LoggerZero{
		logger: logger,
	}
}

func (l *LoggerZero) Debugf(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

func (l *LoggerZero) Infof(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *LoggerZero) Warnf(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *LoggerZero) Errorf(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}
