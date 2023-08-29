# Package terrors

&copy; 2023 Hewlett Packard Enterprise Development LP

```go
import "github.hpe.com/cloud/go-gadgets/x/terrors"
```

Package terrors provides a set of base errors that can be used by applications to reduce boilerplate error code. All errors implement error wrapping and must be used as pointer type (as returned by their constructors).

Whilst the base errors can be used directly within applications, this can lead to a lack of specificity in error handling that can cause problems. It is recommended that instead the base errors are embedded into error types that are specific to the application in question. For example:

	 var VolumeNotFoundType = &VolumeNotFound{}

	 type VolumeNotFound struct {
	     *terrors.NotFound
	 }

	 func NewVolumeNotFound(id uuid.UUID, err error) *VolumeNotFound {
	     return &VolumeNotFound{
	         NotFound: terrors.NewNotFound("volume", "id", id.String()),
		 }
	 }

<details>
<summary><b>Index</b></summary>

- [Variables](#variables)
- [Functions](#functions)
  - [CombineErrsIntoError(msg, errs, infoExtractor)](#func-combineerrsintoerror)
  - [IsPresentable(err)](#func-ispresentable)
  - [Presentable(err)](#func-presentable)
  - [PresentableErrorMsg(err, msg)](#func-presentableerrormsg)
- [Types](#types)
  - [type BaseError](#type-baseerror)
    - [NewBaseError(msg, err)](#func-newbaseerror)
    - [(e) Error()](#func-baseerror-error)
    - [(e) Unwrap()](#func-baseerror-unwrap)
  - [type ErrInfoExtractor](#type-errinfoextractor)
    - [UnwrapInfoExtractor(maxDepth)](#func-unwrapinfoextractor)
  - [type Forbidden](#type-forbidden)
    - [NewForbidden(msg, err)](#func-newforbidden)
    - [(e) Error()](#func-forbidden-error)
    - [(e) Unwrap()](#func-forbidden-unwrap)
  - [type InputAndMsg](#type-inputandmsg)
  - [type InternalError](#type-internalerror)
    - [NewInternalError(msg, err)](#func-newinternalerror)
    - [(e) Error()](#func-internalerror-error)
    - [(e) Unwrap()](#func-internalerror-unwrap)
  - [type InvalidCtxValue](#type-invalidctxvalue)
    - [NewInvalidCtxValue(key)](#func-newinvalidctxvalue)
    - [(e) Error()](#func-invalidctxvalue-error)
    - [(e) Unwrap()](#func-invalidctxvalue-unwrap)
  - [type InvalidInput](#type-invalidinput)
    - [NewInvalidInput(inputs, err)](#func-newinvalidinput)
    - [NewInvalidSingleInput(input, msg, err)](#func-newinvalidsingleinput)
    - [(e) Error()](#func-invalidinput-error)
    - [(e) Unwrap()](#func-invalidinput-unwrap)
  - [type InvalidParameter](#type-invalidparameter)
    - [NewInvalidParameter(param, msg, err)](#func-newinvalidparameter)
    - [NewNilParameter(name)](#func-newnilparameter)
    - [(e) Error()](#func-invalidparameter-error)
    - [(e) Unwrap()](#func-invalidparameter-unwrap)
  - [type NotFound](#type-notfound)
    - [NewNotFound(resource, field, value, err)](#func-newnotfound)
    - [(e) Error()](#func-notfound-error)
    - [(e) Unwrap()](#func-notfound-unwrap)
  - [type NotImplemented](#type-notimplemented)
    - [NewNotImplemented(feature)](#func-newnotimplemented)
    - [(e) Error()](#func-notimplemented-error)
    - [(e) Unwrap()](#func-notimplemented-unwrap)
  - [type PresentableError](#type-presentableerror)
  - [type ServiceUnavailable](#type-serviceunavailable)
    - [NewServiceUnavailable(service, err)](#func-newserviceunavailable)
    - [NewServiceUnavailableWithRetryInterval(service, retryInterval, err)](#func-newserviceunavailablewithretryinterval)
    - [(e) Error()](#func-serviceunavailable-error)
    - [(e) Unwrap()](#func-serviceunavailable-unwrap)
  - [type StateConflict](#type-stateconflict)
    - [NewStateConflict(msg, err)](#func-newstateconflict)
    - [(e) Error()](#func-stateconflict-error)
    - [(e) Unwrap()](#func-stateconflict-unwrap)
  - [type Unauthenticated](#type-unauthenticated)
    - [NewUnauthenticated(msg, err)](#func-newunauthenticated)
    - [(e) Error()](#func-unauthenticated-error)
    - [(e) Unwrap()](#func-unauthenticated-unwrap)
  - [type UnknownError](#type-unknownerror)
    - [NewUnknownError(msg, err)](#func-newunknownerror)
    - [(e) Error()](#func-unknownerror-error)
    - [(e) Unwrap()](#func-unknownerror-unwrap)
  - [type UnsatisfiableRange](#type-unsatisfiablerange)
    - [NewUnsatisfiableRange(firstBytePos, lastBytePos, completeLength, err)](#func-newunsatisfiablerange)
    - [(e) Error()](#func-unsatisfiablerange-error)
    - [(e) Unwrap()](#func-unsatisfiablerange-unwrap)
</details>

<details>
<summary><b>Examples</b></summary>

- [Presentable](#example-presentable)
- [PresentableErrorMsg](#example-presentableerrormsg)

</details>

## Variables

```go
var (
	InternalErrorType	= &InternalError{}
	UnknownErrorType	= &UnknownError{}
	NotImplementedType	= &NotImplemented{}
	InvalidInputType	= &InvalidInput{}
	InvalidParameterType	= &InvalidParameter{}
	NotFoundType		= &NotFound{}
	ServiceUnavailableType	= &ServiceUnavailable{}
	StateConflictType	= &StateConflict{}
	UnauthenticatedType	= &Unauthenticated{}
	ForbiddenType		= &Forbidden{}
	UnsatisfiableRangeType	= &UnsatisfiableRange{}
	InvalidCtxValueType	= &InvalidCtxValue{}
)
```

Error instances to use with the stdlib's errors.As.

These are not thread-safe and should not be used to access the value of a cast error. For example:

	// test whether an error is of the expected type for control flow
	if errors.As(err, &terrors.InternalErrorType) {
	    // do something based on the known type of the error
	}

	// use the value of the error
	internalErr := &terrors.InternalError{}
	if errors.As(err, &internalErr) {
	    // use internalErr to access the error itself
	}


## Functions

### func CombineErrsIntoError

```go
func CombineErrsIntoError(msg string, errs []error, infoExtractor ErrInfoExtractor) error
```

CombineErrsIntoError converts a slice of errors into a single error but combining the error messages. If no extractor function is specified then the default behavior of using err.Error() to get the message is used.


### func IsPresentable

```go
func IsPresentable(err error) bool
```

IsPresentable can be used to check if an error has been marked as presentable. A presentable error that has been wrapped in another error is not considered presentable.


### func Presentable

```go
func Presentable(err error) error
```

Presentable marks an error as 'presentable', meaning that the error message from the err.Error() function is consider safe and appropriate to return to an consumer (e.g. returned in a gRPC or REST response). Presentable error messages should be easy to understand and contain information that the consumer can act upon.

<details>
<summary><b id="example-presentable"><a href="#example-presentable">Example</a></b></summary>



```go
innerErr := errors.New("db connection failed")
internalErr := terrors.Presentable(terrors.NewInternalError("failed to perform DB query", innerErr))

fmt.Println(errors.Is(internalErr, innerErr))

// Output:
// true
```

</details>


### func PresentableErrorMsg

```go
func PresentableErrorMsg(err error, msg string) string
```

PresentableErrorMsg provides a utility to return either the err.Error() message if the error is presentable, or a standard error message if it is not. A presentable error that has been wrapped in another error is not considered presentable.

<details>
<summary><b id="example-presentableerrormsg"><a href="#example-presentableerrormsg">Example</a></b></summary>



```go
errForConsumer := terrors.Presentable(errors.New("page size must be greater than 0"))
wrappedErrForConsumer := fmt.Errorf("pagination token invalid: %w", errForConsumer)
errNotForConsumer := errors.New("db connection failed")

fmt.Println(terrors.PresentableErrorMsg(errForConsumer, "an internal error occurred"))
fmt.Println(terrors.PresentableErrorMsg(wrappedErrForConsumer, "an internal error occurred"))
fmt.Println(terrors.PresentableErrorMsg(errNotForConsumer, "an internal error occurred"))

// Output:
// page size must be greater than 0
// an internal error occurred
// an internal error occurred
```

</details>


## Types

### type BaseError

```go
type BaseError struct {
	// contains filtered or unexported fields
}
```

BaseError is a flexible error for embedding into other error types when the other error types produce messages that would be unnecessary or confusing if they were embedded instead. It is recommended to consider using another error type instead of immediately reaching for this to avoid boilerplate and duplication in errors created using this package.


#### func NewBaseError

```go
func NewBaseError(msg string, err error) *BaseError
```

NewBaseError constructs a new BaseError, wrapping the provided error.


#### func (\*BaseError) Error

```go
func (e *BaseError) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*BaseError) Unwrap

```go
func (e *BaseError) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type ErrInfoExtractor

```go
type ErrInfoExtractor func(error) string
```

ErrInfoExtractor defines a function signature for functions that can extract information from an error.


#### func UnwrapInfoExtractor

```go
func UnwrapInfoExtractor(maxDepth int) ErrInfoExtractor
```

UnwrapInfoExtractor creates an ErrInfoExtractor function that unwraps an error to the specified depth, combining all the messages together into one string.


### type Forbidden

```go
type Forbidden struct {
	Msg string
	// contains filtered or unexported fields
}
```

Forbidden represents a request was received with an acceptable form of authentication, but that the user lacked permission to perform the requested action.


#### func NewForbidden

```go
func NewForbidden(msg string, err error) *Forbidden
```

NewForbidden constructs a Forbidden error. An error can be specified to wrap if it is available.


#### func (\*Forbidden) Error

```go
func (e *Forbidden) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*Forbidden) Unwrap

```go
func (e *Forbidden) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type InputAndMsg

```go
type InputAndMsg struct {
	Input	string
	Msg	string
}
```

InputAndMsg is used to pair an invalid input identifier with a message describing why it is invalid.


### type InternalError

```go
type InternalError struct {
	Msg string
	// contains filtered or unexported fields
}
```

InternalError encompasses logic errors that are not immediately resolvable and that the caller is not expected to perform any actions beyond returning an error itself with a generic error message. For example, the scanning of a SQL row that fails would result in an internal error as either there is a bug in the code or the data has been corrupted.


#### func NewInternalError

```go
func NewInternalError(msg string, err error) *InternalError
```

NewInternalError constructs a new InternalError, wrapping the provided error.


#### func (\*InternalError) Error

```go
func (e *InternalError) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*InternalError) Unwrap

```go
func (e *InternalError) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type InvalidCtxValue

```go
type InvalidCtxValue struct {
	Key string
	// contains filtered or unexported fields
}
```

InvalidCtxValue indicates that a value for the requested key was either not found in the context or was an unexpected type.


#### func NewInvalidCtxValue

```go
func NewInvalidCtxValue(key string) *InvalidCtxValue
```

NewInvalidCtxValue creates a new InvalidCtxValue instance.


#### func (\*InvalidCtxValue) Error

```go
func (e *InvalidCtxValue) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*InvalidCtxValue) Unwrap

```go
func (e *InvalidCtxValue) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type InvalidInput

```go
type InvalidInput struct {
	Inputs []InputAndMsg
	// contains filtered or unexported fields
}
```

InvalidInput represents input from a user/consumer that has been provided that has not passed validation. The inputs are modeled as a slice of input name and message tuples to allow for a single error to capture the validation of a whole input struct. The messages produced from this error are generally expected to suitable for customers to consume and thus should contain any internal information about the service.


#### func NewInvalidInput

```go
func NewInvalidInput(inputs []InputAndMsg, err error) *InvalidInput
```

NewInvalidInput constructs a new InvalidInput error struct. An error can be specified to wrap, but is not expected in most cases.


#### func NewInvalidSingleInput

```go
func NewInvalidSingleInput(input, msg string, err error) *InvalidInput
```

NewInvalidSingleInput constructs a new InvalidInput error struct. An error can be specified to wrap, but is not expected in most cases.


#### func (\*InvalidInput) Error

```go
func (e *InvalidInput) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*InvalidInput) Unwrap

```go
func (e *InvalidInput) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type InvalidParameter

```go
type InvalidParameter struct {
	Parameter	string
	Msg		string
	// contains filtered or unexported fields
}
```

InvalidParameter represents a parameter that does not meet expected criteria. It is differentiated from InvalidInput in that it is expected to be used for defensive programming (e.g. nil checks) rather than for validating external input into the system.


#### func NewInvalidParameter

```go
func NewInvalidParameter(param, msg string, err error) *InvalidParameter
```

NewInvalidParameter constructs a new InvalidParameter error struct. An error can be specified to wrap, but is not expected in most cases.


#### func NewNilParameter

```go
func NewNilParameter(name string) *InvalidParameter
```

NewNilParameter creates a new InvalidParameter instance specialised for reporting unexpectedly nil parameters.


#### func (\*InvalidParameter) Error

```go
func (e *InvalidParameter) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*InvalidParameter) Unwrap

```go
func (e *InvalidParameter) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type NotFound

```go
type NotFound struct {
	Resource	string
	Field		string
	Value		string
	// contains filtered or unexported fields
}
```

NotFound represents that a resource has not been found with the specified identifier. To facilitate interrogation, the resource type is expected to be one of a fixed set of values defined in this package.


#### func NewNotFound

```go
func NewNotFound(resource, field, value string, err error) *NotFound
```

NewNotFound constructs a NotFound error. An error can be specified to wrap if it is available.


#### func (\*NotFound) Error

```go
func (e *NotFound) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*NotFound) Unwrap

```go
func (e *NotFound) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type NotImplemented

```go
type NotImplemented struct {
	Feature string
	// contains filtered or unexported fields
}
```

NotImplemented represents functionality that has been called that has not been implemented yet.


#### func NewNotImplemented

```go
func NewNotImplemented(feature string) *NotImplemented
```

NewNotImplemented constructs a new NotImplemented error. An error can be specified to wrap, but is not expected in most cases.


#### func (\*NotImplemented) Error

```go
func (e *NotImplemented) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*NotImplemented) Unwrap

```go
func (e *NotImplemented) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type PresentableError

```go
type PresentableError interface {
	PresentableError() string
}
```

PresentableError defines an interface that can be implemented to indicate that the message of the error is suitable to be returned as is to the consumer.


### type ServiceUnavailable

```go
type ServiceUnavailable struct {
	Service		string
	RetryInterval	int
	// contains filtered or unexported fields
}
```

ServiceUnavailable represents that an upstream service did not respond or returns an error that indicated it could not complete the request. The latter case may occur if an upstream service to an upstream service did not respond.


#### func NewServiceUnavailable

```go
func NewServiceUnavailable(service string, err error) *ServiceUnavailable
```

NewServiceUnavailable constructs a ServiceUnavailable error. An error can be specified to wrap if it is available.


#### func NewServiceUnavailableWithRetryInterval

```go
func NewServiceUnavailableWithRetryInterval(
	service string,
	retryInterval int,
	err error,
) *ServiceUnavailable
```

NewServiceUnavailableWithRetryInterval constructs a ServiceUnavailable error with the number of seconds to wait until retrying. An error can be specified to wrap if it is available.


#### func (\*ServiceUnavailable) Error

```go
func (e *ServiceUnavailable) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*ServiceUnavailable) Unwrap

```go
func (e *ServiceUnavailable) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type StateConflict

```go
type StateConflict struct {
	Msg string
	// contains filtered or unexported fields
}
```

StateConflict indicates that an operation is valid but not possible due to the current state of the system.


#### func NewStateConflict

```go
func NewStateConflict(msg string, err error) *StateConflict
```

NewStateConflict constructs a StateConflict error. An error can be specified to wrap, but is not expected in most cases.


#### func (\*StateConflict) Error

```go
func (e *StateConflict) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*StateConflict) Unwrap

```go
func (e *StateConflict) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type Unauthenticated

```go
type Unauthenticated struct {
	Msg string
	// contains filtered or unexported fields
}
```

Unauthenticated represents that a request was received without authentication or with an unexpected authenthication scheme, e.g. Basic instead of Bearer.


#### func NewUnauthenticated

```go
func NewUnauthenticated(msg string, err error) *Unauthenticated
```

NewUnauthenticated constructs an Unauthenticated error. An error can be specified to wrap if it is available.


#### func (\*Unauthenticated) Error

```go
func (e *Unauthenticated) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*Unauthenticated) Unwrap

```go
func (e *Unauthenticated) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type UnknownError

```go
type UnknownError struct {
	Msg string
	// contains filtered or unexported fields
}
```

UnknownError encompasses errors that occur when unexpected code paths or similar are reached. For example, when declaring a switch statement that expects one of an explicit set of cases, the default case should return an UnknownError. An UnknownError should almost always be accompanied by an alert in logging.


#### func NewUnknownError

```go
func NewUnknownError(msg string, err error) *UnknownError
```

NewUnknownError constructs a new UnknownError, wrapping the provided error if one is available.


#### func (\*UnknownError) Error

```go
func (e *UnknownError) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*UnknownError) Unwrap

```go
func (e *UnknownError) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


### type UnsatisfiableRange

```go
type UnsatisfiableRange struct {
	FirstBytePos	int64
	LastBytePos	*int64
	CompleteLength	int64
	// contains filtered or unexported fields
}
```

UnsatisfiableRange represents a request that was received with a syntactically valid Range header, but was outside the complete length of the target resource.


#### func NewUnsatisfiableRange

```go
func NewUnsatisfiableRange(
	firstBytePos int64,
	lastBytePos *int64,
	completeLength int64,
	err error,
) *UnsatisfiableRange
```

NewUnsatisfiableRange constructs an UnsatisfiableRange error. An error can be specified to wrap, but is not expected in most cases.


#### func (\*UnsatisfiableRange) Error

```go
func (e *UnsatisfiableRange) Error() string
```

Error allows baseError and any structs that embed it to satisfy the error interface.


#### func (\*UnsatisfiableRange) Unwrap

```go
func (e *UnsatisfiableRange) Unwrap() error
```

Unwrap allows baseError and any structs that embed it to be used with the error wrapping utilities introduced in go 1.13.


