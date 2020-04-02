package hashreader

import (
	stacktrace "github.com/palantir/stacktrace"
)

var errNestedReader = stacktrace.NewError("Nesting of Reader detected, not allowed")

// SHA256Mismatch ...
type SHA256Mismatch struct {
	ExpectedSHA256   string
	CalculatedSHA256 string
}

// Error ...
func (e SHA256Mismatch) Error() string {
	return "Bad sha256: Expected " + e.ExpectedSHA256 + " is not valid with what we calculated " + e.CalculatedSHA256
}

// BadDigest ...
type BadDigest struct {
	ExpectedMD5   string
	CalculatedMD5 string
}

// Error ...
func (e BadDigest) Error() string {
	return "Bad digest: Expected " + e.ExpectedMD5 + " is not valid with what we calculated " + e.CalculatedMD5
}
