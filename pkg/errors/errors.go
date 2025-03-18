// Package errors provides simple error handling utilities.
package errors

import (
	"errors"
	"fmt"
)

// Re-export standard library error functions
var (
	New    = errors.New
	Unwrap = errors.Unwrap
	Is     = errors.Is
	As     = errors.As
)

// Predefined error types
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrIO           = errors.New("io error")
	ErrConfig       = errors.New("config error")
)

// Errorf creates a new error with formatted message
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Errortf creates a new error with the specified type and formatted message
func Errortf(errType error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{errType}, args...)...)
}

// Wrap wraps an existing error with a type and additional context
func Wrapf(err error, errType error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format+": %w", append([]interface{}{errType}, append(args, err)...)...)
}
