package identity

import (
	"context"
	"errors"
	"strings"

	"github.com/workforce-pro/workforce-payroll/internal/foundation"
)

type PrincipalKind string

const (
	PrincipalKindHuman       PrincipalKind = "HUMAN"
	PrincipalKindIntegration PrincipalKind = "INTEGRATION"
)

var (
	ErrUnauthenticated           = errors.New("principal is not authenticated")
	ErrAuthenticationUnavailable = errors.New("authentication service unavailable")
	ErrInvalidClaims             = errors.New("invalid authenticated claims")
)

type VerifiedIdentity struct {
	Issuer  string
	Subject string
	Kind    PrincipalKind
}

func (v VerifiedIdentity) Validate() error {
	if strings.TrimSpace(v.Issuer) == "" || strings.TrimSpace(v.Subject) == "" {
		return ErrInvalidClaims
	}
	if v.Kind != PrincipalKindHuman && v.Kind != PrincipalKindIntegration {
		return ErrInvalidClaims
	}
	return nil
}

type Principal struct {
	ID       foundation.Identifier
	TenantID foundation.Identifier
	Issuer   string
	Subject  string
	Kind     PrincipalKind
}

func (p Principal) Validate() error {
	if p.ID.IsZero() || p.TenantID.IsZero() {
		return ErrUnauthenticated
	}
	return (VerifiedIdentity{Issuer: p.Issuer, Subject: p.Subject, Kind: p.Kind}).Validate()
}

type TokenVerifier interface {
	Verify(ctx context.Context, serialized string) (VerifiedIdentity, error)
}

type PrincipalRepository interface {
	ResolveActive(ctx context.Context, verified VerifiedIdentity) (Principal, error)
}
