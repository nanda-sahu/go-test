// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.hpe.com/cloud/go-gadgets/x/terrors"
)

// Log output keys that are only used in the Zap implementation. These are based on the standardized
// keys listed in https://pages.github.hpe.com/cloud/storage-design/docs/logging.html#log-fields.
const (
	callerKey     = "caller"
	levelKey      = "level"
	messageKey    = "message"
	nameKey       = "name"
	stacktraceKey = "stacktrace"
	timestampKey  = "timestamp"
)

// ZapJSONLogger is an implementation of the logging repository that
// uses zap's sugared logger.
type ZapJSONLogger struct {
	logger *zap.Logger
}

// NewZapJSONLogger creates a zap based logger that implements to Logger repository defined
// in this package. The logger should be flushed before the application exits.
func NewZapJSONLogger(logLevel string, opts ...ZapOption) (*ZapJSONLogger, error) {
	zapConfig, err := stdJSONLoggerConfig(logLevel)
	if err != nil {
		return nil, err
	}

	return zapLoggerFromConfig(zapConfig, opts)
}

// Error logs an error level message. Logs at this level implicitly add a stacktrace field.
func (z *ZapJSONLogger) Error(msg string) {
	z.logger.Error(msg)
}

// Warn logs an warning level message.
func (z *ZapJSONLogger) Warn(msg string) {
	z.logger.Warn(msg)
}

// Info logs an info level message.
func (z *ZapJSONLogger) Info(msg string) {
	z.logger.Info(msg)
}

// Debug logs a debug level message
func (z *ZapJSONLogger) Debug(msg string) {
	z.logger.Debug(msg)
}

// WithFields returns a new logger with the specified key-value pairs attached for
// subsequent logging operations.
// This function returns a repositories logger interface rather than the explicit
// ZapJSONLogger to allow it to satisfy the Logger interface
func (z *ZapJSONLogger) WithFields(fields Fields) Logger {
	fieldList := make([]zap.Field, len(fields))
	i := 0
	for key, value := range fields {
		fieldList[i] = zap.Any(key, value)
		i++
	}

	return &ZapJSONLogger{
		z.logger.With(fieldList...),
	}
}

// WithField returns a new logger with the specified key-value pair attached for
// subsequent logging operations.
// This function returns a repositories logger interface rather than the explicit
// ZapJSONLogger to allow it to satisfy the Logger interface
func (z *ZapJSONLogger) WithField(key string, value interface{}) Logger {
	return &ZapJSONLogger{
		z.logger.With(zap.Any(key, value)),
	}
}

// WithError provides a wrapper around WithField to add an error field to the logger,
// ensuring consistency of error message keys. It will also unwrap the error, unlike a
// normal WithField call.
func (z *ZapJSONLogger) WithError(err error) Logger {
	unwrapper := terrors.UnwrapInfoExtractor(1000) //nolint:gomnd // arbitrary exit condition to avoid infinite loop
	msg := unwrapper(err)
	return z.WithField(ErrKey, msg)
}

// Flush flushes any pending log statements. This is a no-op as logs are written to STDOUT and
// synchonization is not supported on STDOUT/STDERR.
func (z *ZapJSONLogger) Flush() error {
	return nil
}

func zapLoggerFromConfig(config zap.Config, opts []ZapOption) (*ZapJSONLogger, error) {
	zapOpts := defaultZapOptions()

	for _, opt := range opts {
		opt.apply(zapOpts)
	}

	logger, err := config.Build(
		zap.AddCallerSkip(zapOpts.callerSkipCount),
		zap.Hooks(zapOpts.hooks...),
	)
	if err != nil {
		return nil, err
	}

	return &ZapJSONLogger{
		logger: logger,
	}, nil
}

func stdJSONLoggerConfig(logLevel string) (zap.Config, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return zap.Config{}, err
	}

	return zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        timestampKey,
			LevelKey:       levelKey,
			MessageKey:     messageKey,
			NameKey:        nameKey,
			StacktraceKey:  stacktraceKey,
			CallerKey:      callerKey,
			LineEnding:     "\n",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}, nil
}
