package foundation

import (
	"errors"
	"strings"
	"time"
)

// CommandEnvelope carries the full trust context for any mutation command.
// Every write operation in the system is wrapped in a CommandEnvelope before
// reaching domain logic. The server populates Principal and Tenant from the
// authenticated context; callers supply idempotency, versioning, and
// correlation metadata.
//
// Requirements: 1.7–1.10, 7.1–7.10
type CommandEnvelope struct {
	// PrincipalID is the server-resolved authenticated actor.
	PrincipalID Identifier

	// TenantID is the server-derived effective Tenant; never caller-supplied.
	TenantID Identifier

	// PrincipalKind distinguishes HUMAN from INTEGRATION principals.
	PrincipalKind string

	// Permission is the required permission code for the action.
	Permission string

	// ResourceScope is the Tenant sub-scope (org unit, site, team, worker, etc.)
	// required by the action, if any.
	ResourceScope string

	// CorrelationID traces a logical business action across components.
	CorrelationID string

	// CausationID identifies the immediate cause (parent command or event).
	CausationID string

	// IdempotencyKey is a Tenant-scoped unique client key for safe retries.
	// Requirements: 7.1–7.4
	IdempotencyKey string

	// ExpectedVersion is the caller's last-known aggregate version for
	// optimistic concurrency control.
	// Requirements: 7.5–7.7
	ExpectedVersion *int64

	// RequestDigest is the canonical SHA-256 digest of the submitted payload,
	// used to detect idempotency-key reuse with different content.
	// Requirements: 7.8–7.9
	RequestDigest Digest

	// IssuedAt records when the command was accepted by the server boundary.
	IssuedAt time.Time
}

// Validate ensures the envelope has all mandatory server-supplied fields.
func (cmd CommandEnvelope) Validate() error {
	if cmd.PrincipalID.IsZero() {
		return errors.New("principal_id is required")
	}
	if cmd.TenantID.IsZero() {
		return errors.New("tenant_id is required")
	}
	if strings.TrimSpace(cmd.PrincipalKind) == "" {
		return errors.New("principal_kind is required")
	}
	if strings.TrimSpace(cmd.CorrelationID) == "" {
		return errors.New("correlation_id is required")
	}
	if strings.TrimSpace(cmd.IdempotencyKey) == "" {
		return errors.New("idempotency_key is required")
	}
	if cmd.RequestDigest.IsZero() {
		return errors.New("request_digest is required")
	}
	if cmd.IssuedAt.IsZero() {
		return errors.New("issued_at is required")
	}
	return nil
}

// CheckVersion compares the expected version against the current aggregate
// version. Returns nil if no expected version was supplied (caller accepts
// any version) or if versions match. Returns ErrConflict otherwise.
func (cmd CommandEnvelope) CheckVersion(currentVersion int64) error {
	if cmd.ExpectedVersion == nil {
		return nil
	}
	if *cmd.ExpectedVersion != currentVersion {
		return ErrConflict
	}
	return nil
}
