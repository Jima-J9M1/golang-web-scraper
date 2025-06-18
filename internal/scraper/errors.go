package scraper

import "fmt"

type ErrFetchFailed struct {
	URL        string
	StatusCode int
	Reason     string
	WrappedErr error
}

func (e *ErrFetchFailed) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("failed to fetch url %s (status %d) Reason %s error %w", e.URL, e.StatusCode, e.Reason, e.WrappedErr)
	}

	return fmt.Sprintf("failed to fetch URL %s: %s: %v", e.URL, e.Reason, e.WrappedErr)
}

func (e *ErrFetchFailed) Unwrap() error {
	return e.WrappedErr
}

type ErrParseFailed struct {
	URL string
	Reason string
	WrappedErr error
}


func (e *ErrParseFailed) Error() string {
	return fmt.Sprintf("failed to parse HTML from %s: %s: %v", e.URL, e.Reason, e.WrappedErr)
}

// Unwrap allows errors.Is and errors.As to inspect the wrapped error.
func (e *ErrParseFailed) Unwrap() error {
	return e.WrappedErr
}

// ErrValidation represents an error due to invalid input (e.g., malformed URL).
type ErrValidation struct {
	Field      string
	Value      string
	Constraint string
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("validation failed for field '%s' (value '%s'): %s", e.Field, e.Value, e.Constraint)
}