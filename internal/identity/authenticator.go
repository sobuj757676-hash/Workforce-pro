package identity

import (
	"context"
	"errors"
)

type Authenticator struct {
	verifier   TokenVerifier
	principals PrincipalRepository
}

func NewAuthenticator(verifier TokenVerifier, principals PrincipalRepository) (*Authenticator, error) {
	if verifier == nil || principals == nil {
		return nil, errors.New("token verifier and principal repository are required")
	}
	return &Authenticator{verifier: verifier, principals: principals}, nil
}

func (a *Authenticator) Authenticate(ctx context.Context, serializedToken string) (Principal, error) {
	verified, err := a.verifier.Verify(ctx, serializedToken)
	if err != nil {
		return Principal{}, ErrUnauthenticated
	}
	principal, err := a.principals.ResolveActive(ctx, verified)
	if errors.Is(err, ErrPrincipalNotFound) {
		return Principal{}, ErrUnauthenticated
	}
	if err != nil {
		return Principal{}, ErrAuthenticationUnavailable
	}
	if err := principal.Validate(); err != nil {
		return Principal{}, ErrUnauthenticated
	}
	return principal, nil
}
