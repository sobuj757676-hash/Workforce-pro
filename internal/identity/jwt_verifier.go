package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var ErrInvalidToken = errors.New("invalid bearer token")

const principalKindClaim = "principal_kind"

type JWTVerifier struct {
	issuer         string
	audience       string
	keys           jwk.Set
	acceptableSkew time.Duration
}

func NewJWTVerifier(issuer, audience string, keys jwk.Set, acceptableSkew time.Duration) (*JWTVerifier, error) {
	if strings.TrimSpace(issuer) == "" || strings.TrimSpace(audience) == "" || keys == nil || keys.Len() == 0 {
		return nil, errors.New("issuer, audience, and verification keys are required")
	}
	if acceptableSkew < 0 {
		return nil, errors.New("acceptable skew cannot be negative")
	}
	return &JWTVerifier{issuer: issuer, audience: audience, keys: keys, acceptableSkew: acceptableSkew}, nil
}

func (v *JWTVerifier) Verify(ctx context.Context, serialized string) (VerifiedIdentity, error) {
	if strings.TrimSpace(serialized) == "" {
		return VerifiedIdentity{}, ErrInvalidToken
	}

	token, err := jwt.ParseString(serialized,
		jwt.WithKeySet(v.keys),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithAcceptableSkew(v.acceptableSkew),
		jwt.WithRequiredClaim(jwt.ExpirationKey),
		jwt.WithRequiredClaim(jwt.IssuedAtKey),
		jwt.WithRequiredClaim(jwt.IssuerKey),
		jwt.WithRequiredClaim(jwt.SubjectKey),
		jwt.WithRequiredClaim(jwt.AudienceKey),
		jwt.WithRequiredClaim(principalKindClaim),
		jwt.WithContext(ctx),
	)
	if err != nil {
		return VerifiedIdentity{}, ErrInvalidToken
	}

	issuer, issuerOK := token.Issuer()
	subject, subjectOK := token.Subject()
	var kindValue string
	if !issuerOK || !subjectOK || token.Get(principalKindClaim, &kindValue) != nil {
		return VerifiedIdentity{}, ErrInvalidToken
	}

	verified := VerifiedIdentity{Issuer: issuer, Subject: subject, Kind: PrincipalKind(kindValue)}
	if err := verified.Validate(); err != nil {
		return VerifiedIdentity{}, ErrInvalidToken
	}
	return verified, nil
}

func ParseBearerAuthorization(header string) (string, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", ErrInvalidToken
	}
	return parts[1], nil
}
