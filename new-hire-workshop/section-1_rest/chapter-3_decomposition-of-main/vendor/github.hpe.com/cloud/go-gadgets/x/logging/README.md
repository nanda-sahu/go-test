# Package logging

&copy; 2023 Hewlett Packard Enterprise Development LP

```go
import "github.hpe.com/cloud/go-gadgets/x/logging"
```

Package logging provides a standard logging interface for applications to use. It also includes a Zap implementation of that interface for use.

In general, logs should be sufficient to diagnose issues / determine if the system is running as desired, but not so verbose as to detract from their usage. Superfluous logs, such as at the start/end of every function, should not be present as this is likely to overwhelm log ingestion services, e.g. Humio.

4 levels of logging are supported in line with the levels specified in the logging standards which should be used according to the below guidelines:

  - Debug: Information that is desired only for debugging bench/non-production workloads. Usage of this log level should be low, since it is only for debugging during development and not diagnosing production.
  - Info: Events that are useful to know about (e.g. for debugging a live production workload), but are not an issue. Using info level logs for events in the happy path allows investigation into whether the application is behaving in the expected manner (by inspection of volume/frequency of given info level logs).
  - Warn: Any event that may be an error, but from which the application is currently able to continue. Likely to be used infrequently as generally it is immediately apparent whether an event is an error or not. A potential use case may be logging 4xx client errors for a REST API which can usually be ignored unless the volume is exceedingly high.
  - Error: Any event that is not in the happy path. It is expected that if an unhappy path that starts from deeper within the application, e.g. a failed database connection, the process would likely have multiple error logs as it propagates back up the system. If clean architecture is used, then an example of these multiple error logs be when the error is observed in the datastore and use case layers. The log should include enough information to investigate a production issue that resulted from this event.

It is often tempting to embed a value into a log message. However, this is almost always better served by adding a field to the logger:

	// don't do this
	logger.Error(fmt.Sprintf("could not find ID %s for example resource: %v", id, err))

	// do this instead
	logger.WithError(err).WithField("id", id).Error("could not find example resource")

By leveraging fields properly and so making the most of structured logging, the values can be filtered to in a log service by using a key-value match rather than trying to match the whole log message. This is particularly useful if a field appears multiple times.

<details>
<summary><b id="example-package-fields"><a href="#example-package-fields">Example (Fields)</a></b></summary>

Fields can be attached to the log as key-value pairs which appear as separate fields once injested into a log service.


```go
// add key-value pairs to the logs one at a time
logger.WithField("key", "value").Info("example")

// or multiple pairs in one go
logger.
	WithFields(logging.Fields{
		"key1":	"value1",
		"key2":	"value2",
		"key3":	"value3",
	}).
	Info("example")

// there's a special helper for adding errors and unwrapping them
logger.WithError(err).Error("an unexpected error occurred")

// and you can add fields to the logger to be used later
logger = logger.WithField("key", "value")
logger.Info("message with a key:value pair attached")
```

</details>

<details>
<summary><b id="example-package-levels"><a href="#example-package-levels">Example (Levels)</a></b></summary>

Multiple log levels are supported that can be used to signify the success or failure of an operation.


```go
logger.Error("this is an error")
logger.Warn("this is a warning")
logger.Info("this is for information")
logger.Debug("this is for debugging")
```

</details>

<details>
<summary><b>Index</b></summary>

- [Constants](#constants)
- [Functions](#functions)
  - [ContextWithLogger(ctx, logger)](#func-contextwithlogger)
- [Types](#types)
  - [type Fields](#type-fields)
  - [type Logger](#type-logger)
    - [LoggerFromContext(ctx)](#func-loggerfromcontext)
    - [NewNoopLogger()](#func-newnooplogger)
  - [type NoopLogger](#type-nooplogger)
    - [() Debug()](#func-nooplogger-debug)
    - [() Error()](#func-nooplogger-error)
    - [() Flush()](#func-nooplogger-flush)
    - [() Info()](#func-nooplogger-info)
    - [() Warn()](#func-nooplogger-warn)
    - [(l) WithError()](#func-nooplogger-witherror)
    - [(l) WithField()](#func-nooplogger-withfield)
    - [(l) WithFields()](#func-nooplogger-withfields)
  - [type ZapJSONLogger](#type-zapjsonlogger)
    - [NewZapJSONLogger(logLevel, opts)](#func-newzapjsonlogger)
    - [(z) Debug(msg)](#func-zapjsonlogger-debug)
    - [(z) Error(msg)](#func-zapjsonlogger-error)
    - [(z) Flush()](#func-zapjsonlogger-flush)
    - [(z) Info(msg)](#func-zapjsonlogger-info)
    - [(z) Warn(msg)](#func-zapjsonlogger-warn)
    - [(z) WithError(err)](#func-zapjsonlogger-witherror)
    - [(z) WithField(key, value)](#func-zapjsonlogger-withfield)
    - [(z) WithFields(fields)](#func-zapjsonlogger-withfields)
  - [type ZapOption](#type-zapoption)
    - [WithHooks(hooks)](#func-withhooks)
    - [WithSkipCallerCount(count)](#func-withskipcallercount)
</details>

<details>
<summary><b>Examples</b></summary>

- [Package (Fields)](#example-package-fields)
- [Package (Levels)](#example-package-levels)

</details>

## Constants

```go
const (
	CustomerIDKey	= "customer-id"
	ErrKey		= "error"
	PanicKey	= "panic"
	SpanIDKey	= "x-b3-spanid"
	TraceIDKey	= "x-b3-traceid"
	UserIDKey	= "user-id"
)
```

Log field keys for adding standardized fields to log messages. These are based on those listed in [https://pages.github.hpe.com/cloud/storage-design/docs/logging.html#log-fields](https://pages.github.hpe.com/cloud/storage-design/docs/logging.html#log-fields).


## Functions

### func ContextWithLogger

```go
func ContextWithLogger(ctx context.Context, logger Logger) context.Context
```

ContextWithLogger returns a copy of the given context with the given logger added as a value. To extract the stored logger, use LoggerFromContext.


## Types

### type Fields

```go
type Fields map[string]interface{}
```



### type Logger

```go
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
```

Logger defines the repository interface for a generic application logger, intended to result in structured logging using the WithField(s) functions to add context to the logger.

Implementations are expected to provide new instances of the Logger when returning from the WithField(s) functions to allow for the creation of child loggers that's subsequent use don't influence the parent.


#### func LoggerFromContext

```go
func LoggerFromContext(ctx context.Context) (Logger, error)
```

LoggerFromContext extracts a logger instance from the given context. If the logger was not found or was an unexpected type, an InvalidCtxValue error will be returned. To add a logger to the context, use ContextWithLogger.


#### func NewNoopLogger

```go
func NewNoopLogger() Logger
```

NewNoopLogger creates a new NoopLogger instance.


### type NoopLogger

```go
type NoopLogger struct{}
```

NoopLogger is a mock implementation of the logging.Logger interface, intended to be used specifically in contract tests, where log assertions is not required.


#### func (\*NoopLogger) Debug

```go
func (*NoopLogger) Debug(string)
```

Debug ignored by logger.


#### func (\*NoopLogger) Error

```go
func (*NoopLogger) Error(string)
```

Error ignored by logger.


#### func (\*NoopLogger) Flush

```go
func (*NoopLogger) Flush() error
```

Flush flushes any pending log statements. This is a no-op as no logs are stored.


#### func (\*NoopLogger) Info

```go
func (*NoopLogger) Info(string)
```

Info ignored by logger.


#### func (\*NoopLogger) Warn

```go
func (*NoopLogger) Warn(string)
```

Warn ignored by logger.


#### func (\*NoopLogger) WithError

```go
func (l *NoopLogger) WithError(error) Logger
```

WithError ignored by logger.


#### func (\*NoopLogger) WithField

```go
func (l *NoopLogger) WithField(string, interface{}) Logger
```

WithField ignored by logger.


#### func (\*NoopLogger) WithFields

```go
func (l *NoopLogger) WithFields(Fields) Logger
```

WithFields ignored by logger.


### type ZapJSONLogger

```go
type ZapJSONLogger struct {
	// contains filtered or unexported fields
}
```

ZapJSONLogger is an implementation of the logging repository that uses zap's sugared logger.


#### func NewZapJSONLogger

```go
func NewZapJSONLogger(logLevel string, opts ...ZapOption) (*ZapJSONLogger, error)
```

NewZapJSONLogger creates a zap based logger that implements to Logger repository defined in this package. The logger should be flushed before the application exits.


#### func (\*ZapJSONLogger) Debug

```go
func (z *ZapJSONLogger) Debug(msg string)
```

Debug logs a debug level message


#### func (\*ZapJSONLogger) Error

```go
func (z *ZapJSONLogger) Error(msg string)
```

Error logs an error level message. Logs at this level implicitly add a stacktrace field.


#### func (\*ZapJSONLogger) Flush

```go
func (z *ZapJSONLogger) Flush() error
```

Flush flushes any pending log statements. This is a no-op as logs are written to STDOUT and synchonization is not supported on STDOUT/STDERR.


#### func (\*ZapJSONLogger) Info

```go
func (z *ZapJSONLogger) Info(msg string)
```

Info logs an info level message.


#### func (\*ZapJSONLogger) Warn

```go
func (z *ZapJSONLogger) Warn(msg string)
```

Warn logs an warning level message.


#### func (\*ZapJSONLogger) WithError

```go
func (z *ZapJSONLogger) WithError(err error) Logger
```

WithError provides a wrapper around WithField to add an error field to the logger, ensuring consistency of error message keys. It will also unwrap the error, unlike a normal WithField call.


#### func (\*ZapJSONLogger) WithField

```go
func (z *ZapJSONLogger) WithField(key string, value interface{}) Logger
```

WithField returns a new logger with the specified key-value pair attached for subsequent logging operations. This function returns a repositories logger interface rather than the explicit ZapJSONLogger to allow it to satisfy the Logger interface


#### func (\*ZapJSONLogger) WithFields

```go
func (z *ZapJSONLogger) WithFields(fields Fields) Logger
```

WithFields returns a new logger with the specified key-value pairs attached for subsequent logging operations. This function returns a repositories logger interface rather than the explicit ZapJSONLogger to allow it to satisfy the Logger interface


### type ZapOption

```go
type ZapOption interface {
	// contains filtered or unexported methods
}
```

ZapOption configures a ZapJSONLogger


#### func WithHooks

```go
func WithHooks(hooks ...func(zapcore.Entry) error) ZapOption
```

WithHooks creates a ZapOption that adds hooks to the internal zap.Logger


#### func WithSkipCallerCount

```go
func WithSkipCallerCount(count int) ZapOption
```

WithSkipCallerCount creates a ZapOption that overrides the default caller skip count.

This is helpful for skipping wrappers of the logger when printing the caller name. Default value is 1.


