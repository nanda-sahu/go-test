# Changelog

This is the changelog for the `x/logging` module. For instructions on how to add to this, see the
[contributing guide](../../CONTRIBUTING.md).

## Legend

Changes are grouped to describe their impact on the module, as follows:

- `Added` for new APIs or fields.
- `Changed` for changes in existing functionality.
- `Deprecated` for identifying which previously stable features will be removed.
- `Removed` for listing deprecated or unstable features that were removed.
- `Fixed` for any bug fixes.
- `Security` for changes made to resolve vulnerabilities.

## Unreleased

## [0.0.4] - 2023-02-06

### Added

- Added `NoopLogger` to provide an ability to ignore reported logs, when they are not required to be asserted.
- Added example usages of the Logger interface.
- Added context accessors that can return an `terrors.InvalidCtxValue` error.
    - `logging.LoggerFromContext` for accessing the logger. A logger can also be set using `logging.ContextWithLogger` in test scenarios.
- Added `assertlogging.LoggerBuilder` to support defining a logger in table-driven tests.
- Add type-safe methods in generated mocks.

### Changed

- Improved debuggability of unexpected or unfulfilled log assertions by including the log message and fields in the error message.
- Upgrade Go to 1.19.
- `NewZapJSONLogger` changed to except a list of options for improved flexibility. Usage: `NewZapJSONLogger(logLevel, WithHooks(hooks), WithSkipCallerCount(callerCount))`
- Updated x/terrors to v0.0.3.

## [0.0.3] - 2022-09-21

### Changed

- Upgraded to Go 1.17.
- Make `assertlogging.Logger` threadsafe so it can be used by multiple goroutines in unit tests.
- Stub out `ZapJSONLogger.Flush` as synchronization of `STDOUT` is not supported.
- Updated x/terrors to v0.0.2.

### Added

- Documented the requirement for ERROR logs to include a stacktrace.
- Added `assertlogging.AssertWithoutOrder` as an option for `assertlogging.NewLogger` to allow logs to be matched out of
  order. This is intended to be used when a single logger is consumed by multiple goroutines and so could produce logs
  in a variety of orders.

## [0.0.2] - 2022-05-17

### Changed

- Added `assertlogging/NewField` to make asserting against multiple log fields using `assertlogging.Logger.WithFields` more ergonomic.

### Removed

- Removed the `CallerKey`, `LevelKey`, `MessageKey`, `MethodKey`, `ServiceKey`, and `TimestampKey` constants as they
  were either considered internal details of the `ZapJSONLogger` or unnecessary.

### Added

- Added a `PanicKey` constant for logging values passed to `panic` and retrieved via `recover`.
- Added a range of new helpers for log field assertions:
    - `Equalf`: Similar to `Equal`, but uses string formatting to allow for partially dynamic log field values.
    - `EqualErrorf`: Similar to `EqualError`, but uses string formatting to allow for partially dynamic log field values.
    - `True`: For checking whether a log field value is `true`.
    - `False`: For checking whether a log field value is `false`.
    - `Empty`: For checking for empty log field values in strings, slices, etc.
    - `Nil`: For checking whether a log field value is `nil`.

## [0.0.1] - 2022-03-01

### Added

- Add a `Logger` interface to implement the [DSCC Logging standards](https://pages.github.hpe.com/cloud/storage-design/docs/logging.html).
- Add `ZapJSONLogger` to implement the `Logger` interface using [`go.uber.org/zap`](https://go.uber.org/zap).
