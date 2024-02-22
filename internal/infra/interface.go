package infra

// Logger defines the logging interface.
type Logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

// Notifier defines the notification interface.
type Notifier interface {
	Sendf(topic, msg string)
}
