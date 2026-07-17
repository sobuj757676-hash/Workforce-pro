package httpapi

import (
	"errors"
	"net/http"

	"github.com/workforce-pro/workforce-payroll/internal/identity"
	"github.com/workforce-pro/workforce-payroll/internal/tenancy"
)

func Authenticate(authenticator *identity.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			serialized, err := identity.ParseBearerAuthorization(request.Header.Get("Authorization"))
			if err != nil {
				writeUnauthenticated(response)
				return
			}
			principal, err := authenticator.Authenticate(request.Context(), serialized)
			if errors.Is(err, identity.ErrAuthenticationUnavailable) {
				writeAuthenticationUnavailable(response)
				return
			}
			if err != nil {
				writeUnauthenticated(response)
				return
			}
			ctx, err := tenancy.WithPrincipal(request.Context(), principal)
			if err != nil {
				writeUnauthenticated(response)
				return
			}
			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}

func EnforceSuppliedTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		effective, err := tenancy.EffectiveTenant(request.Context())
		if err != nil {
			writeUnauthenticated(response)
			return
		}

		supplied := append([]string(nil), request.Header.Values("X-Tenant-ID")...)
		supplied = append(supplied, request.URL.Query()["tenant_id"]...)
		for _, candidate := range supplied {
			if err := tenancy.CompareSuppliedTenant(effective, candidate); err != nil {
				writeTenantScopedNotFound(response)
				return
			}
		}
		next.ServeHTTP(response, request)
	})
}
