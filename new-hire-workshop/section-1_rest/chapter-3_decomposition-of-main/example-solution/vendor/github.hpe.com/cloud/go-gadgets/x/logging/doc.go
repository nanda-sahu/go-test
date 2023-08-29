// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

// Package logging provides a standard logging interface for applications to use. It also
// includes a Zap implementation of that interface for use.
//
// In general, logs should be sufficient to diagnose issues / determine if the system is
// running as desired, but not so verbose as to detract from their usage. Superfluous logs,
// such as at the start/end of every function, should not be present as this is likely to overwhelm
// log ingestion services, e.g. Humio.
//
// 4 levels of logging are supported in line with the levels specified in the logging standards
// which should be used according to the below guidelines:
//
//   - Debug: Information that is desired only for debugging bench/non-production workloads.
//     Usage of this log level should be low, since it is only for debugging during development and
//     not diagnosing production.
//   - Info: Events that are useful to know about (e.g. for debugging a live production
//     workload), but are not an issue. Using info level logs for events in the happy path allows
//     investigation into whether the application is behaving in the expected manner (by inspection of
//     volume/frequency of given info level logs).
//   - Warn: Any event that may be an error, but from which the application is currently able to
//     continue. Likely to be used infrequently as generally it is immediately apparent whether
//     an event is an error or not. A potential use case may be logging 4xx client errors for a REST
//     API which can usually be ignored unless the volume is exceedingly high.
//   - Error: Any event that is not in the happy path. It is expected that if an unhappy path that
//     starts from deeper within the application, e.g. a failed database connection, the process
//     would likely have multiple error logs as it propagates back up the system. If clean
//     architecture is used, then an example of these multiple error logs be when the error is
//     observed in the datastore and use case layers. The log should include enough information
//     to investigate a production issue that resulted from this event.
//
// It is often tempting to embed a value into a log message. However, this is almost always better
// served by adding a field to the logger:
//
//	// don't do this
//	logger.Error(fmt.Sprintf("could not find ID %s for example resource: %v", id, err))
//
//	// do this instead
//	logger.WithError(err).WithField("id", id).Error("could not find example resource")
//
// By leveraging fields properly and so making the most of structured logging, the values can be
// filtered to in a log service by using a key-value match rather than trying to match the whole
// log message. This is particularly useful if a field appears multiple times.
package logging
