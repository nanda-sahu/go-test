// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package logging

// Non-allocating compile time check to ensure the logging.Logger interface is implemented
// correctly.
var _ Logger = &NoopLogger{}

// NoopLogger is a mock implementation of the logging.Logger interface, intended to be used
// specifically in contract tests, where log assertions is not required.
type NoopLogger struct{}

// NewNoopLogger creates a new NoopLogger instance.
func NewNoopLogger() Logger {
	return &NoopLogger{}
}

// Error ignored by logger.
func (*NoopLogger) Error(string) {}

// Warn ignored by logger.
func (*NoopLogger) Warn(string) {}

// Info ignored by logger.
func (*NoopLogger) Info(string) {}

// Debug ignored by logger.
func (*NoopLogger) Debug(string) {}

// WithField ignored by logger.
func (l *NoopLogger) WithField(string, interface{}) Logger {
	return l
}

// WithError ignored by logger.
func (l *NoopLogger) WithError(error) Logger {
	return l
}

// WithFields ignored by logger.
func (l *NoopLogger) WithFields(Fields) Logger {
	return l
}

// Flush flushes any pending log statements. This is a no-op as no logs are stored.
func (*NoopLogger) Flush() error {
	return nil
}
