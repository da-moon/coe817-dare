package errors

import "github.com/palantir/stacktrace"

// Errors
var (
	// ErrInvalidPayloadSize ...
	ErrInvalidPayloadSize = stacktrace.NewError("invalid payload size")
	// ErrAuthentication ...
	ErrAuthentication = stacktrace.NewError("authentication failed")
	// ErrNonceMismatch ...
	ErrNonceMismatch = stacktrace.NewError("header nonce mismatch")
	// ErrUnexpectedEOF ...
	ErrUnexpectedEOF = stacktrace.NewError("unexpected end of file (EOF)")
	// ErrUnexpectedData ...
	ErrUnexpectedData = stacktrace.NewError("unexpected data after final burst of data")
)
