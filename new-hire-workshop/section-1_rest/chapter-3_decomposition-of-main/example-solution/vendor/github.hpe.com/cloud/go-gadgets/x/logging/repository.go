// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package logging

// Log field keys for adding standardized fields to log messages. These are based on those listed in
// https://pages.github.hpe.com/cloud/storage-design/docs/logging.html#log-fields.
const (
	CustomerIDKey = "customer-id"
	ErrKey        = "error"
	PanicKey      = "panic"
	SpanIDKey     = "x-b3-spanid"
	TraceIDKey    = "x-b3-traceid"
	UserIDKey     = "user-id"
)

type Fields map[string]interface{}

// Logger defines the repository interface for a generic application logger,
// intended to result in structured logging using the WithField(s) functions to
// add context to the logger.
//
// Implementations are expected to provide new instances of the Logger when returning from
// the WithField(s) functions to allow for the creation of child loggers that's subsequent use
// don't influence the parent.
type Logger interface {
	// Used when an error has occurred that is not recoverable, and will most likely
	// involve returning an error to the consumer/user. Implementations must include a stacktrace at this level.
	Error(msg string)

	// Used when a potential issue may exist, but the system can continue to function.
	Warn(msg string)

	// Used when something of interest has occurred that is useful to have logged in a
	// production setting.
	Info(msg string)

	// Used when providing information on specific code paths with the application that are
	// being executed that are not required in a production setting.
	Debug(msg string)

	// WithField returns a new instance of the Logger that has the specified field attached
	// in all subsequent messages.
	WithField(key string, value interface{}) Logger

	// WithError provides a wrapper around WithField to add an error field to the logger,
	// ensuring consistency of error message keys.
	WithError(err error) Logger

	// WithFields returns a new instance of the Logger that has the specified fields attached
	// in all subsequent messages.
	WithFields(fields Fields) Logger

	// Flush ensures that any pending log messages are written out. For some implementations
	// this function will be a no-op.
	Flush() error
}
