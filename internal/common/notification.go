package common

// Logger defines the logging interface.
type Notifier interface {
	Sendf(topic, msg string)
}
