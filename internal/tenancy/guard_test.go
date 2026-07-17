package tenancy

import (
	"context"
	"errors"
	"testing"

	"github.com/workforce-pro/workforce-payroll/internal/foundation"
	"github.com/workforce-pro/workforce-payroll/internal/identity"
	"pgregory.net/rapid"
)

const (
	tenantOne    = "018f6f4a-7b2c-7d11-8a33-123456789abc"
	tenantTwo    = "018f6f4a-7b2c-7d11-9b44-abcdef123456"
	principalOne = "018f6f4a-7b2c-7d11-aa55-001122334455"
	resourceOne  = "018f6f4a-7b2c-7d11-bb66-556677889900"
)

func authenticatedContext(t *testing.T, tenant string) context.Context {
	t.Helper()
	ctx, err := WithPrincipal(context.Background(), identity.Principal{
		ID:       foundation.MustParseIdentifier(principalOne),
		TenantID: foundation.MustParseIdentifier(tenant),
		Issuer:   "https://identity.example.test/realms/workforce",
		Subject:  "subject",
		Kind:     identity.PrincipalKindHuman,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestCompareSuppliedTenant(t *testing.T) {
	effective := foundation.MustParseIdentifier(tenantOne)
	if err := CompareSuppliedTenant(effective, ""); err != nil {
		t.Fatalf("omitted Tenant rejected: %v", err)
	}
	if err := CompareSuppliedTenant(effective, tenantOne); err != nil {
		t.Fatalf("matching Tenant rejected: %v", err)
	}
	for _, supplied := range []string{tenantTwo, "invalid"} {
		if !errors.Is(CompareSuppliedTenant(effective, supplied), ErrTenantScopedNotFound) {
			t.Fatalf("supplied Tenant %q did not receive non-disclosing denial", supplied)
		}
	}
}

type recordedResource struct {
	tenant foundation.Identifier
}

func (resource recordedResource) ResourceTenantID() foundation.Identifier { return resource.tenant }

type scopedReader struct {
	resource        recordedResource
	found           bool
	requestedTenant foundation.Identifier
}

func (reader *scopedReader) FindInTenant(_ context.Context, tenantID, _ foundation.Identifier) (recordedResource, error) {
	reader.requestedTenant = tenantID
	if !reader.found || reader.resource.tenant != tenantID {
		return recordedResource{}, ErrTenantScopedNotFound
	}
	return reader.resource, nil
}

func TestLoadResourceDoesNotDiscloseCrossTenantExistence(t *testing.T) {
	reader := &scopedReader{resource: recordedResource{tenant: foundation.MustParseIdentifier(tenantTwo)}, found: true}
	_, err := LoadResource(authenticatedContext(t, tenantOne), reader, foundation.MustParseIdentifier(resourceOne))
	if !errors.Is(err, ErrTenantScopedNotFound) {
		t.Fatalf("LoadResource() error = %v", err)
	}
	if reader.requestedTenant != foundation.MustParseIdentifier(tenantOne) {
		t.Fatalf("repository received Tenant %s, want effective Tenant %s", reader.requestedTenant, tenantOne)
	}
}

// Feature: workforce-payroll-management-saas, Property 1: Tenant Isolation.
// For every caller-supplied Tenant different from the server-derived Tenant,
// authorization denies without returning resource data or mutating state.
// **Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.6**
func TestProperty1TenantIsolation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		effective := generatedIdentifier(t, "effective Tenant")
		other := generatedIdentifier(t, "resource Tenant")
		if other == effective {
			other[15] ^= 1
		}
		principalID := generatedIdentifier(t, "Principal")
		resourceID := generatedIdentifier(t, "resource")
		kind := rapid.SampledFrom([]identity.PrincipalKind{
			identity.PrincipalKindHuman,
			identity.PrincipalKindIntegration,
		}).Draw(t, "Principal kind")
		ctx, err := WithPrincipal(context.Background(), identity.Principal{
			ID: principalID, TenantID: effective, Issuer: "https://identity.example.test", Subject: "generated", Kind: kind,
		})
		if err != nil {
			t.Fatal(err)
		}

		authoritativeState := rapid.Int64().Draw(t, "initial state")
		before := authoritativeState
		reader := &scopedReader{resource: recordedResource{tenant: other}, found: true}
		mutate := func(_ context.Context, _ recordedResource) error {
			authoritativeState++
			return nil
		}

		// A resource owned by another Tenant is invisible even when the caller
		// supplies its own valid Tenant value.
		err = ExecuteMutation(ctx, effective.String(), reader, resourceID, mutate)
		if !errors.Is(err, ErrTenantScopedNotFound) {
			t.Fatalf("cross-Tenant resource returned error %v", err)
		}
		if reader.requestedTenant != effective {
			t.Fatalf("repository received Tenant %s, want %s", reader.requestedTenant, effective)
		}

		// Caller-supplied substitution is rejected before repository or mutation.
		reader.requestedTenant = foundation.Identifier{}
		err = ExecuteMutation(ctx, other.String(), reader, resourceID, mutate)
		if !errors.Is(err, ErrTenantScopedNotFound) {
			t.Fatalf("Tenant substitution returned error %v", err)
		}
		if !reader.requestedTenant.IsZero() {
			t.Fatalf("repository was called after Tenant substitution with %s", reader.requestedTenant)
		}
		if authoritativeState != before {
			t.Fatalf("cross-Tenant denial changed state from %d to %d", before, authoritativeState)
		}
	})
}

func generatedIdentifier(t *rapid.T, label string) foundation.Identifier {
	bytes := rapid.SliceOfN(rapid.Byte(), 16, 16).Draw(t, label)
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	var id foundation.Identifier
	copy(id[:], bytes)
	return id
}
