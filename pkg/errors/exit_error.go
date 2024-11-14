package errors

import (
	"io"

	"github.com/fastly/cli/pkg/text"
)

// SkipExitError is an error that can cause the os.Exit(1) to be skipped.
// An example is 'help' output (e.g. --help).
type SkipExitError struct {
	Skip bool
	Err  error
}

// Unwrap returns the inner error.
func (ee SkipExitError) Unwrap() error {
	return ee.Err
}

// Error prints the inner error string.
func (ee SkipExitError) Error() string {
	if ee.Err == nil {
		return ""
	}
	return ee.Err.Error()
}

// Print the error to the io.Writer for human consumption.
// The inner error is always printed via text.Output with an "Error: " prefix
// and a "." suffix.
func (ee SkipExitError) Print(w io.Writer) {
	if ee.Err != nil {
		text.Error(w, "%s.", ee.Err.Error())
	}
}

// PassthroughExitError is an error that carries an 'exit code' from a
// subprocess, which should be returned from this program. An example
// is `fastly run`.
type PassthroughExitError struct {
	ExitCode int
	Err      error
}

// Unwrap returns the inner error.
func (pte PassthroughExitError) Unwrap() error {
	return pte.Err
}

// Error prints the inner error string.
func (pte PassthroughExitError) Error() string {
	if pte.Err == nil {
		return ""
	}
	return pte.Err.Error()
}
