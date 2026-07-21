package foundation

import (
	"context"
	"time"
)

// TenantScopedEntity is the minimal interface for any aggregate or entity
// that participates in Tenant-scoped storage. Every persisted domain record
// must carry and expose its owning Tenant.
// Requirements: 1.2–1.6
type TenantScopedEntity interface {
	ScopedID() TenantScopedID
}

// VersionedEntity extends TenantScopedEntity with optimistic concurrency.
// Requirements: 7.1–7.10
type VersionedEntity interface {
	TenantScopedEntity
	Version() int64
	ContentDigest() Digest
}

// IdempotencyRecord stores the proof of a previously accepted mutation
// so that safe retries converge on one equivalent server effect.
// Requirements: 7.1–7.4
type IdempotencyRecord struct {
	TenantID      Identifier
	Key           string
	RequestDigest Digest
	ResultRef     string
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

// IdempotencyStore provides the storage contract for idempotency keys.
// The store is always Tenant-scoped.
type IdempotencyStore interface {
	// Find returns an existing idempotency record or ErrNotFound.
	Find(ctx context.Context, tenantID Identifier, key string) (IdempotencyRecord, error)

	// Store persists a new idempotency record. If a record with the same
	// Tenant+key already exists, it returns ErrConflict.
	Store(ctx context.Context, record IdempotencyRecord) error
}

// ObjectKeyBuilder constructs Tenant-scoped keys for object storage.
// Requirements: 1.5 — object keys are Tenant scoped.
type ObjectKeyBuilder struct {
	Prefix string
}

// Build produces a Tenant-scoped object key in the form:
// {prefix}/{tenant_id}/{entity_id}
func (b ObjectKeyBuilder) Build(id TenantScopedID) string {
	return b.Prefix + "/" + id.TenantID.String() + "/" + id.EntityID.String()
}

// CacheKeyBuilder constructs Tenant-scoped cache keys.
// Requirements: 1.5 — cache entries are Tenant scoped.
type CacheKeyBuilder struct {
	Namespace string
}

// Build produces a Tenant-scoped cache key in the form:
// {namespace}:{tenant_id}:{entity_id}
func (b CacheKeyBuilder) Build(id TenantScopedID) string {
	return b.Namespace + ":" + id.TenantID.String() + ":" + id.EntityID.String()
}

// MessageEnvelope wraps any domain event or message with mandatory
// Tenant-scoped metadata for queue/event delivery.
// Requirements: 1.5 — messages are Tenant scoped.
type MessageEnvelope struct {
	TenantID      Identifier
	EventID       Identifier
	AggregateID   TenantScopedID
	CorrelationID string
	CausationID   string
	OccurredAt    time.Time
	SchemaVersion string
	EventType     string
}

// ExportReference identifies a Tenant-scoped export artifact.
// Requirements: 1.5 — exports are Tenant scoped.
type ExportReference struct {
	TenantID    Identifier
	ExportID    Identifier
	Period      string
	GeneratedAt time.Time
	Version     string
}
