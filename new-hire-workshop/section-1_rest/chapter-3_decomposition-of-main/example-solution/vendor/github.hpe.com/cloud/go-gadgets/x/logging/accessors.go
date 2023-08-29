// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package logging

import (
	"context"
	"reflect"

	"github.hpe.com/cloud/go-gadgets/x/terrors"
)

// ctxKey is used in place of a string when adding values to a context to avoid type-based
// collisions.
type ctxKey string

const (
	loggerCtxKey ctxKey = "logger"
)

// LoggerFromContext extracts a logger instance from the given context. If the logger was not found
// or was an unexpected type, an InvalidCtxValue error will be returned. To add a logger to the
// context, use ContextWithLogger.
func LoggerFromContext(ctx context.Context) (Logger, error) {
	logger, ok := ctx.Value(loggerCtxKey).(Logger)
	if !ok || reflect.ValueOf(logger).IsNil() {
		return nil, terrors.NewInvalidCtxValue(string(loggerCtxKey))
	}

	return logger, nil
}

// ContextWithLogger returns a copy of the given context with the given logger added as a value. To
// extract the stored logger, use LoggerFromContext.
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
