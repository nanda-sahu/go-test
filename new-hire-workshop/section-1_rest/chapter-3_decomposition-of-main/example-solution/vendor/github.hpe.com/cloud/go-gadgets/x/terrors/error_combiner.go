// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package terrors

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInfoExtractor defines a function signature for functions
// that can extract information from an error.
type ErrInfoExtractor func(error) string

// UnwrapInfoExtractor creates an ErrInfoExtractor function that unwraps an error
// to the specified depth, combining all the messages together into one string.
func UnwrapInfoExtractor(maxDepth int) ErrInfoExtractor {
	return func(err error) string {
		if err == nil {
			return ""
		}

		builder := strings.Builder{}
		builder.WriteString(err.Error())

		for i := 1; i < maxDepth; i++ {
			err = errors.Unwrap(err)
			if err == nil {
				return builder.String()
			}

			builder.WriteString(": " + err.Error())
		}

		return builder.String()
	}
}

// CombineErrsIntoError converts a slice of errors into a single error but combining the
// error messages. If no extractor function is specified then the default behavior of
// using err.Error() to get the message is used.
func CombineErrsIntoError(msg string, errs []error, infoExtractor ErrInfoExtractor) error {
	if len(errs) == 0 {
		return nil
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s: ", msg))
	for i, err := range errs {
		if i != 0 {
			builder.WriteString(", ")
		}

		errMsg := err.Error()
		if infoExtractor != nil {
			errMsg = infoExtractor(err)
		}
		builder.WriteString(errMsg)
	}

	return errors.New(builder.String())
}
