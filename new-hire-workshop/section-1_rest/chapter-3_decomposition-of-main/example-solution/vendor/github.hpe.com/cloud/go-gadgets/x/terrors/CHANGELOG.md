# Changelog

This is the changelog for the `x/terrors` module. For instructions on how to add to this, see the
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

## [0.0.3] - 2023-02-03

### Added

- Added `NilParameter` for when a function argument was unexpectedly nil.
- Add type-safe methods in generated mocks.

### Changed

- Upgrade Go to 1.19.

## [0.0.2] - 2022-09-20

### Added

- Added constructors for generated mocks.
- Added `BaseError` for embedding into new error types where embedding a different error type might produce an unecessary or confusing error message.
- Added `InvalidCtxValue` for when a context accessor could not find a correctly typed value within a context.

### Changed

- Upgrade Go to 1.17.

### Fixed

- Fixed a panic when `nil` was passed to `UnwrapInfoExtractor` inline with how `errors.Unwrap` behaves.

## [0.0.1] - 2022-01-25

### Added

- A series of common error types that can be embedded to create your own error types or used directly.
- A `PresentableError` interface for providing error messages that are considered appropriate for API consumers.
