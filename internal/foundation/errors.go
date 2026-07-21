package foundation

import "errors"

// Sentinel errors shared across the domain layer.
// These are stable codes that map to safe HTTP responses.
var (
	// ErrNotFound means the resource does not exist within the effective Tenant scope.
	ErrNotFound = errors.New("resource not found")

	// ErrConflict means an optimistic concurrency or idempotency check failed.
	ErrConflict = errors.New("version conflict")

	// ErrForbidden means the Principal lacks the required Permission or Resource_Scope.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidInput means the request payload failed validation before command execution.
	ErrInvalidInput = errors.New("invalid input")

	// ErrDigestMismatch means the submitted content digest does not match the current version.
	ErrDigestMismatch = errors.New("content digest mismatch")

	// ErrIdempotencyKeyReuse means the same idempotency key was reused with a different request digest.
	ErrIdempotencyKeyReuse = errors.New("idempotency key reused with different request")
)
