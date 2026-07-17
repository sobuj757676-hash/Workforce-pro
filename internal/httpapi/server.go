package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/workforce-pro/workforce-payroll/internal/identity"
	"github.com/workforce-pro/workforce-payroll/internal/tenancy"
)

type Server struct {
	router http.Handler
}

func NewServer(authenticator *identity.Authenticator) (*Server, error) {
	if authenticator == nil {
		return nil, identity.ErrUnauthenticated
	}

	router := chi.NewRouter()
	router.Use(Authenticate(authenticator))
	router.Use(EnforceSuppliedTenant)
	router.Get("/v1/context", getContext)
	router.Get("/v1/tenants/{tenantID}/context", getContextForSuppliedTenant)
	return &Server{router: router}, nil
}

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(response, request)
}

type contextResponse struct {
	PrincipalID string                 `json:"principal_id"`
	TenantID    string                 `json:"tenant_id"`
	Kind        identity.PrincipalKind `json:"principal_kind"`
}

func getContext(response http.ResponseWriter, request *http.Request) {
	principal, err := tenancy.PrincipalFromContext(request.Context())
	if err != nil {
		writeUnauthenticated(response)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(response).Encode(contextResponse{
		PrincipalID: principal.ID.String(),
		TenantID:    principal.TenantID.String(),
		Kind:        principal.Kind,
	})
}

func getContextForSuppliedTenant(response http.ResponseWriter, request *http.Request) {
	effective, err := tenancy.EffectiveTenant(request.Context())
	if err != nil {
		writeUnauthenticated(response)
		return
	}
	if err := tenancy.CompareSuppliedTenant(effective, chi.URLParam(request, "tenantID")); err != nil {
		writeTenantScopedNotFound(response)
		return
	}
	getContext(response, request)
}
