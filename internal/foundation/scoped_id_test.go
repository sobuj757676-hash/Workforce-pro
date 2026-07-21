package foundation

import "testing"

func TestTenantScopedIDValidation(t *testing.T) {
	tenantID := MustParseIdentifier("018f6f4a-7b2c-7d11-8a33-123456789abc")
	entityID := MustParseIdentifier("018f6f4a-7b2c-7d11-bb66-556677889900")

	id, err := NewTenantScopedID(tenantID, entityID)
	if err != nil {
		t.Fatalf("NewTenantScopedID() error = %v", err)
	}
	if id.TenantID != tenantID || id.EntityID != entityID {
		t.Fatalf("unexpected ID components: %+v", id)
	}
	if id.IsZero() {
		t.Fatal("valid TenantScopedID reported as zero")
	}
	if !id.BelongsTo(tenantID) {
		t.Fatal("BelongsTo() returned false for owning Tenant")
	}
	otherTenant := MustParseIdentifier("018f6f4a-7b2c-7d11-9b44-abcdef123456")
	if id.BelongsTo(otherTenant) {
		t.Fatal("BelongsTo() returned true for non-owning Tenant")
	}
}

func TestTenantScopedIDRejectsZeroComponents(t *testing.T) {
	valid := MustParseIdentifier("018f6f4a-7b2c-7d11-8a33-123456789abc")
	if _, err := NewTenantScopedID(Identifier{}, valid); err == nil {
		t.Fatal("accepted zero tenant")
	}
	if _, err := NewTenantScopedID(valid, Identifier{}); err == nil {
		t.Fatal("accepted zero entity")
	}
}

func TestTenantScopedIDString(t *testing.T) {
	tenantID := MustParseIdentifier("018f6f4a-7b2c-7d11-8a33-123456789abc")
	entityID := MustParseIdentifier("018f6f4a-7b2c-7d11-bb66-556677889900")
	id, _ := NewTenantScopedID(tenantID, entityID)
	want := "018f6f4a-7b2c-7d11-8a33-123456789abc/018f6f4a-7b2c-7d11-bb66-556677889900"
	if got := id.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
