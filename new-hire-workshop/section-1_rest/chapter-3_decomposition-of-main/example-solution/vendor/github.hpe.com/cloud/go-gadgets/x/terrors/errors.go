// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package terrors

import (
	"fmt"
	"strconv"
	"strings"
)

// Error instances to use with the stdlib's errors.As.
//
// These are not thread-safe and should not be used to access the value of a cast error. For
// example:
//
//	// test whether an error is of the expected type for control flow
//	if errors.As(err, &terrors.InternalErrorType) {
//	    // do something based on the known type of the error
//	}
//
//	// use the value of the error
//	internalErr := &terrors.InternalError{}
//	if errors.As(err, &internalErr) {
//	    // use internalErr to access the error itself
//	}
var (
	InternalErrorType      = &InternalError{}
	UnknownErrorType       = &UnknownError{}
	NotImplementedType     = &NotImplemented{}
	InvalidInputType       = &InvalidInput{}
	InvalidParameterType   = &InvalidParameter{}
	NotFoundType           = &NotFound{}
	ServiceUnavailableType = &ServiceUnavailable{}
	StateConflictType      = &StateConflict{}
	UnauthenticatedType    = &Unauthenticated{}
	ForbiddenType          = &Forbidden{}
	UnsatisfiableRangeType = &UnsatisfiableRange{}
	InvalidCtxValueType    = &InvalidCtxValue{}
)

// baseError represents a generic error from the domain package that provides
// 'error' functionality to the rest of the typed errors in the package.
type baseError struct {
	Err error
	msg string
}

// Error allows baseError and any structs that embed it to satisfy the error
// interface.
func (e *baseError) Error() string {
	return e.msg
}

// Unwrap allows baseError and any structs that embed it to be used with the
// error wrapping utilities introduced in go 1.13.
func (e *baseError) Unwrap() error {
	// This nil check accounts for the situation where the embedded *baseError
	// in one of the public errors is nil - if it has been constructed without
	// using one of the helper functions (e.g in other package's unit tests).
	if e == nil {
		return nil
	}
	return e.Err
}

// newBaseError is a constructor for a base error. It should not be called
// directly outside of constructing other errors in this package.
//
// NOTE: This function deliberately returns a value rather than a pointer
// as baseError needs to be embedded as a value into the other errors to
// ensure that only the pointer receivors of those errors satisfy the error
// interface.
func newBaseError(msg string, err error) baseError {
	return baseError{
		Err: err,
		msg: msg,
	}
}

// BaseError is a flexible error for embedding into other error types when the
// other error types produce messages that would be unnecessary or confusing if
// they were embedded instead. It is recommended to consider using another error
// type instead of immediately reaching for this to avoid boilerplate and
// duplication in errors created using this package.
type BaseError struct {
	baseError
}

// NewBaseError constructs a new BaseError, wrapping the provided error.
func NewBaseError(msg string, err error) *BaseError {
	return &BaseError{
		baseError: newBaseError(msg, err),
	}
}

// InternalError encompasses logic errors that are not immediately resolvable
// and that the caller is not expected to perform any actions beyond returning
// an error itself with a generic error message. For example, the scanning of a
// SQL row that fails would result in an internal error as either there is a bug
// in the code or the data has been corrupted.
type InternalError struct {
	baseError
	Msg string
}

// NewInternalError constructs a new InternalError, wrapping the provided error.
func NewInternalError(msg string, err error) *InternalError {
	return &InternalError{
		baseError: newBaseError(
			fmt.Sprintf("an internal error occurred: %s", msg),
			err,
		),
		Msg: msg,
	}
}

// UnknownError encompasses errors that occur when unexpected code paths or
// similar are reached. For example, when declaring a switch statement that
// expects one of an explicit set of cases, the default case should return an
// UnknownError. An UnknownError should almost always be accompanied by an alert
// in logging.
type UnknownError struct {
	baseError
	Msg string
}

// NewUnknownError constructs a new UnknownError, wrapping the provided error if
// one is available.
func NewUnknownError(msg string, err error) *UnknownError {
	return &UnknownError{
		baseError: newBaseError(
			fmt.Sprintf("an unexpected error occurred: %s", msg),
			err,
		),
		Msg: msg,
	}
}

// NotImplemented represents functionality that has been called that has not
// been implemented yet.
type NotImplemented struct {
	baseError
	Feature string
}

// NewNotImplemented constructs a new NotImplemented error. An error can be
// specified to wrap, but is not expected in most cases.
func NewNotImplemented(feature string) *NotImplemented {
	return &NotImplemented{
		baseError: newBaseError(
			fmt.Sprintf("feature %s is not implemented", feature),
			nil,
		),
		Feature: feature,
	}
}

// InputAndMsg is used to pair an invalid input identifier with a message
// describing why it is invalid.
type InputAndMsg struct {
	Input string
	Msg   string
}

// InvalidInput represents input from a user/consumer that has been provided
// that has not passed validation. The inputs are modeled as a slice of input
// name and message tuples to allow for a single error to capture the validation
// of a whole input struct. The messages produced from this error are generally
// expected to suitable for customers to consume and thus should contain any
// internal information about the service.
type InvalidInput struct {
	baseError
	Inputs []InputAndMsg
}

// NewInvalidInput constructs a new InvalidInput error struct. An error can be
// specified to wrap, but is not expected in most cases.
func NewInvalidInput(inputs []InputAndMsg, err error) *InvalidInput {
	builder := strings.Builder{}
	for i, inputAndMsg := range inputs {
		if i != 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%s %s", inputAndMsg.Input, inputAndMsg.Msg))
	}
	return &InvalidInput{
		baseError: newBaseError(
			fmt.Sprintf("invalid input(s): %s", builder.String()),
			err,
		),
		Inputs: inputs,
	}
}

// NewInvalidSingleInput constructs a new InvalidInput error struct. An error
// can be specified to wrap, but is not expected in most cases.
func NewInvalidSingleInput(input, msg string, err error) *InvalidInput {
	return NewInvalidInput([]InputAndMsg{
		{
			Input: input,
			Msg:   msg,
		},
	}, err)
}

// InvalidParameter represents a parameter that does not meet expected criteria.
// It is differentiated from InvalidInput in that it is expected to be used for
// defensive programming (e.g. nil checks) rather than for validating external
// input into the system.
type InvalidParameter struct {
	baseError
	Parameter string
	Msg       string
}

// NewInvalidParameter constructs a new InvalidParameter error struct. An error
// can be specified to wrap, but is not expected in most cases.
func NewInvalidParameter(param, msg string, err error) *InvalidParameter {
	return &InvalidParameter{
		baseError: newBaseError(
			fmt.Sprintf("invalid parameter %s: %s", param, msg),
			err,
		),
		Parameter: param,
		Msg:       msg,
	}
}

// NewNilParameter creates a new InvalidParameter instance specialised for reporting unexpectedly
// nil parameters.
func NewNilParameter(name string) *InvalidParameter {
	return NewInvalidParameter(name, "must not be nil", nil)
}

// NotFound represents that a resource has not been found with the specified
// identifier. To facilitate interrogation, the resource type is expected to be
// one of a fixed set of values defined in this package.
type NotFound struct {
	baseError
	Resource string
	Field    string
	Value    string
}

// NewNotFound constructs a NotFound error. An error can be specified to wrap if
// it is available.
func NewNotFound(resource, field, value string, err error) *NotFound {
	return &NotFound{
		baseError: newBaseError(
			fmt.Sprintf("%s was not found with %s: %s", resource, field, value),
			err,
		),
		Resource: resource,
		Field:    field,
		Value:    value,
	}
}

// ServiceUnavailable represents that an upstream service did not respond or
// returns an error that indicated it could not complete the request. The latter
// case may occur if an upstream service to an upstream service did not respond.
type ServiceUnavailable struct {
	baseError
	Service       string
	RetryInterval int
}

// NewServiceUnavailable constructs a ServiceUnavailable error. An error can be
// specified to wrap if it is available.
func NewServiceUnavailable(service string, err error) *ServiceUnavailable {
	return &ServiceUnavailable{
		baseError: newBaseError(
			fmt.Sprintf("%s service unavailable", service),
			err,
		),
		Service: service,
	}
}

// NewServiceUnavailableWithRetryInterval constructs a ServiceUnavailable error
// with the number of seconds to wait until retrying. An error can be specified
// to wrap if it is available.
func NewServiceUnavailableWithRetryInterval(
	service string,
	retryInterval int,
	err error,
) *ServiceUnavailable {
	return &ServiceUnavailable{
		baseError: newBaseError(
			fmt.Sprintf("%s service unavailable", service),
			err,
		),
		Service:       service,
		RetryInterval: retryInterval,
	}
}

// StateConflict indicates that an operation is valid but not possible due to
// the current state of the system.
type StateConflict struct {
	baseError
	Msg string
}

// NewStateConflict constructs a StateConflict error. An error can be specified
// to wrap, but is not expected in most cases.
func NewStateConflict(msg string, err error) *StateConflict {
	return &StateConflict{
		baseError: newBaseError(
			fmt.Sprintf("operation not possible due to current state: %s", msg),
			err,
		),
		Msg: msg,
	}
}

// Unauthenticated represents that a request was received without authentication
// or with an unexpected authenthication scheme, e.g. Basic instead of Bearer.
type Unauthenticated struct {
	baseError
	Msg string
}

// NewUnauthenticated constructs an Unauthenticated error. An error can be
// specified to wrap if it is available.
func NewUnauthenticated(msg string, err error) *Unauthenticated {
	return &Unauthenticated{
		baseError: newBaseError(
			fmt.Sprintf("unauthenticated: %s", msg),
			err,
		),
		Msg: msg,
	}
}

// Forbidden represents a request was received with an acceptable form of
// authentication, but that the user lacked permission to perform the requested
// action.
type Forbidden struct {
	baseError
	Msg string
}

// NewForbidden constructs a Forbidden error. An error can be specified to wrap
// if it is available.
func NewForbidden(msg string, err error) *Forbidden {
	return &Forbidden{
		baseError: newBaseError(
			fmt.Sprintf("forbidden: %s", msg),
			err,
		),
		Msg: msg,
	}
}

// UnsatisfiableRange represents a request that was received with a
// syntactically valid Range header, but was outside the complete length of the
// target resource.
type UnsatisfiableRange struct {
	baseError
	FirstBytePos   int64
	LastBytePos    *int64
	CompleteLength int64
}

// NewUnsatisfiableRange constructs an UnsatisfiableRange error. An error can be
// specified to wrap, but is not expected in most cases.
func NewUnsatisfiableRange(
	firstBytePos int64,
	lastBytePos *int64,
	completeLength int64,
	err error,
) *UnsatisfiableRange {
	byteRange := fmt.Sprintf("%d-", firstBytePos)
	if lastBytePos != nil {
		byteRange += strconv.FormatInt(*lastBytePos, 10) //nolint:gomnd // numbers are in base 10
	}

	return &UnsatisfiableRange{
		baseError: newBaseError(
			fmt.Sprintf(
				"unsatisfiable range: %s (complete length: %d)",
				byteRange,
				completeLength,
			),
			err,
		),
		FirstBytePos:   firstBytePos,
		LastBytePos:    lastBytePos,
		CompleteLength: completeLength,
	}
}

// InvalidCtxValue indicates that a value for the requested key was either not found in the context
// or was an unexpected type.
type InvalidCtxValue struct {
	baseError
	Key string
}

// NewInvalidCtxValue creates a new InvalidCtxValue instance.
func NewInvalidCtxValue(key string) *InvalidCtxValue {
	return &InvalidCtxValue{
		baseError: newBaseError(
			fmt.Sprintf("context value %s was missing or not the expected type", key),
			nil,
		),
		Key: key,
	}
}
