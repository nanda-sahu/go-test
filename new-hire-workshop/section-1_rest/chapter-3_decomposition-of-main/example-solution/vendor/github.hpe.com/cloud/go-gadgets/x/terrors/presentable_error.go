// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package terrors

// PresentableError defines an interface that can be implemented to indicate that
// the message of the error is suitable to be returned as is to the consumer.
type PresentableError interface {
	PresentableError() string
}

type presentableWrapper struct {
	err error
}

func (c *presentableWrapper) Error() string {
	return c.err.Error()
}

func (c *presentableWrapper) Unwrap() error {
	return c.err
}

func (c *presentableWrapper) PresentableError() string {
	return c.Error()
}

// Presentable marks an error as 'presentable', meaning that the error message
// from the err.Error() function is consider safe and appropriate to return
// to an consumer (e.g. returned in a gRPC or REST response).
// Presentable error messages should be easy to understand and contain information
// that the consumer can act upon.
func Presentable(err error) error {
	return &presentableWrapper{
		err: err,
	}
}

// IsPresentable can be used to check if an error has been marked as presentable.
// A presentable error that has been wrapped in another error is not considered presentable.
func IsPresentable(err error) bool {
	_, ok := err.(PresentableError)
	return ok
}

// PresentableErrorMsg provides a utility to return either the err.Error() message
// if the error is presentable, or a standard error message if it is not.
// A presentable error that has been wrapped in another error is not considered presentable.
func PresentableErrorMsg(err error, msg string) string {
	e, ok := err.(PresentableError)
	if !ok {
		return msg
	}
	return e.PresentableError()
}
