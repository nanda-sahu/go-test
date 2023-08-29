// (C) Copyright 2021-2023 Hewlett Packard Enterprise Development LP

package logging

import (
	"go.uber.org/zap/zapcore"
)

type zapOptions struct {
	callerSkipCount int
	hooks           []func(zapcore.Entry) error
}

func defaultZapOptions() *zapOptions {
	return &zapOptions{
		callerSkipCount: 1, // ignore the logger itself when providing caller
		hooks:           []func(zapcore.Entry) error{},
	}
}

// ZapOption configures a ZapJSONLogger
type ZapOption interface {
	apply(*zapOptions)
}

// WithHooks creates a ZapOption that adds hooks to the internal zap.Logger
func WithHooks(hooks ...func(zapcore.Entry) error) ZapOption {
	return hooksOption{hooks}
}

type hooksOption struct {
	hooks []func(zapcore.Entry) error
}

func (h hooksOption) apply(options *zapOptions) {
	options.hooks = h.hooks
}

// WithSkipCallerCount creates a ZapOption that overrides the default caller skip count.
//
// This is helpful for skipping wrappers of the logger when printing the caller name.
// Default value is 1.
func WithSkipCallerCount(count int) ZapOption {
	return callerCountOption{count}
}

type callerCountOption struct {
	count int
}

func (c callerCountOption) apply(options *zapOptions) {
	options.callerSkipCount = c.count
}
