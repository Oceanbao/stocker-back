package common

// Logger defines the logging interface.
type Logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
