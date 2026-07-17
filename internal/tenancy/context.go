package tenancy

import (
	"context"
	"errors"

	"github.com/workforce-pro/workforce-payroll/internal/foundation"
	"github.com/workforce-pro/workforce-payroll/internal/identity"
)

var ErrPrincipalContextMissing = errors.New("authenticated principal context is missing")

type principalContextKey struct{}

func WithPrincipal(ctx context.Context, principal identity.Principal) (context.Context, error) {
	if err := principal.Validate(); err != nil {
		return nil, ErrPrincipalContextMissing
	}
	return context.WithValue(ctx, principalContextKey{}, principal), nil
}

func PrincipalFromContext(ctx context.Context) (identity.Principal, error) {
	principal, ok := ctx.Value(principalContextKey{}).(identity.Principal)
	if !ok || principal.Validate() != nil {
		return identity.Principal{}, ErrPrincipalContextMissing
	}
	return principal, nil
}

func EffectiveTenant(ctx context.Context) (foundation.Identifier, error) {
	principal, err := PrincipalFromContext(ctx)
	if err != nil {
		return foundation.Identifier{}, err
	}
	return principal.TenantID, nil
}
