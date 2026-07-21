package foundation

import "errors"

// TenantScopedID is a canonical pair that binds any entity identifier
// to exactly one Tenant. It is used at trust boundaries to ensure that
// all domain identifiers carry their ownership context.
// Requirements: 1.2–1.6 — every record, cache, object key, message,
// idempotency record, export, and event is Tenant scoped.
type TenantScopedID struct {
	TenantID Identifier
	EntityID Identifier
}

// NewTenantScopedID constructs a TenantScopedID and validates that neither
// component is zero.
func NewTenantScopedID(tenantID, entityID Identifier) (TenantScopedID, error) {
	id := TenantScopedID{TenantID: tenantID, EntityID: entityID}
	if err := id.Validate(); err != nil {
		return TenantScopedID{}, err
	}
	return id, nil
}

// Validate ensures both identifiers are non-zero RFC 4122 UUIDs.
func (id TenantScopedID) Validate() error {
	if id.TenantID.IsZero() {
		return errors.New("tenant identifier is required")
	}
	if id.EntityID.IsZero() {
		return errors.New("entity identifier is required")
	}
	return nil
}

// IsZero returns true when either component is the zero value.
func (id TenantScopedID) IsZero() bool {
	return id.TenantID.IsZero() || id.EntityID.IsZero()
}

// String returns a slash-separated canonical representation.
func (id TenantScopedID) String() string {
	return id.TenantID.String() + "/" + id.EntityID.String()
}

// BelongsTo returns true when the resource is owned by the given Tenant.
func (id TenantScopedID) BelongsTo(tenantID Identifier) bool {
	return id.TenantID == tenantID
}
