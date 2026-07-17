package tenancy

import (
	"context"
	"errors"

	"github.com/workforce-pro/workforce-payroll/internal/foundation"
)

var ErrTenantScopedNotFound = errors.New("tenant-scoped resource not found")

// CompareSuppliedTenant treats mismatch and malformed caller input as the same
// non-disclosing outcome. An omitted Tenant is valid because the server context
// remains authoritative.
func CompareSuppliedTenant(effective foundation.Identifier, supplied string) error {
	if supplied == "" {
		return nil
	}
	candidate, err := foundation.ParseIdentifier(supplied)
	if err != nil || candidate != effective {
		return ErrTenantScopedNotFound
	}
	return nil
}

type ScopedResource interface {
	ResourceTenantID() foundation.Identifier
}

type ScopedResourceReader[T ScopedResource] interface {
	FindInTenant(ctx context.Context, tenantID, resourceID foundation.Identifier) (T, error)
}

type ScopedMutation[T ScopedResource] func(ctx context.Context, resource T) error

// ExecuteMutation performs both caller-supplied and resource ownership checks
// before invoking the authoritative mutation. A boundary failure therefore
// cannot call the mutation function.
func ExecuteMutation[T ScopedResource](
	ctx context.Context,
	suppliedTenant string,
	reader ScopedResourceReader[T],
	resourceID foundation.Identifier,
	mutate ScopedMutation[T],
) error {
	effectiveTenant, err := EffectiveTenant(ctx)
	if err != nil {
		return ErrTenantScopedNotFound
	}
	if err := CompareSuppliedTenant(effectiveTenant, suppliedTenant); err != nil {
		return ErrTenantScopedNotFound
	}
	resource, err := LoadResource(ctx, reader, resourceID)
	if err != nil {
		return ErrTenantScopedNotFound
	}
	if mutate == nil {
		return errors.New("mutation function is required")
	}
	return mutate(ctx, resource)
}

// LoadResource queries only inside the effective Tenant. A repository must not
// perform a global ID lookup before this call because that would create an
// existence side channel.
func LoadResource[T ScopedResource](
	ctx context.Context,
	reader ScopedResourceReader[T],
	resourceID foundation.Identifier,
) (T, error) {
	var zero T
	effectiveTenant, err := EffectiveTenant(ctx)
	if err != nil {
		return zero, ErrTenantScopedNotFound
	}
	resource, err := reader.FindInTenant(ctx, effectiveTenant, resourceID)
	if err != nil {
		return zero, ErrTenantScopedNotFound
	}
	if resource.ResourceTenantID() != effectiveTenant {
		return zero, ErrTenantScopedNotFound
	}
	return resource, nil
}
