package identity

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const (
	testIssuer   = "https://identity.example.test/realms/workforce"
	testAudience = "workforce-api"
)

type jwtFixture struct {
	privateKey jwk.Key
	publicKey  jwk.Key
	keys       jwk.Set
	verifier   *JWTVerifier
}

func newJWTFixture(t *testing.T) jwtFixture {
	t.Helper()
	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := jwk.Import(rawKey)
	if err != nil {
		t.Fatal(err)
	}
	if err := privateKey.Set("kid", "test-key"); err != nil {
		t.Fatal(err)
	}
	if err := privateKey.Set("alg", jwa.RS256()); err != nil {
		t.Fatal(err)
	}
	publicKey, err := privateKey.PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	keys := jwk.NewSet()
	if err := keys.AddKey(publicKey); err != nil {
		t.Fatal(err)
	}
	verifier, err := NewJWTVerifier(testIssuer, testAudience, keys, 0)
	if err != nil {
		t.Fatal(err)
	}
	return jwtFixture{privateKey: privateKey, publicKey: publicKey, keys: keys, verifier: verifier}
}

func (fixture jwtFixture) sign(t *testing.T, kind string, mutate func(*jwt.Builder)) string {
	t.Helper()
	now := time.Now().UTC()
	builder := jwt.NewBuilder().
		Issuer(testIssuer).
		Subject("principal-subject").
		Audience([]string{testAudience}).
		IssuedAt(now.Add(-time.Second)).
		Expiration(now.Add(time.Minute)).
		Claim(principalKindClaim, kind)
	if mutate != nil {
		mutate(builder)
	}
	token, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), fixture.privateKey))
	if err != nil {
		t.Fatal(err)
	}
	return string(serialized)
}

func TestJWTVerifierAcceptsHumanAndIntegrationPrincipals(t *testing.T) {
	fixture := newJWTFixture(t)
	for _, kind := range []PrincipalKind{PrincipalKindHuman, PrincipalKindIntegration} {
		t.Run(string(kind), func(t *testing.T) {
			verified, err := fixture.verifier.Verify(context.Background(), fixture.sign(t, string(kind), nil))
			if err != nil {
				t.Fatalf("Verify() error = %v", err)
			}
			if verified.Kind != kind || verified.Issuer != testIssuer || verified.Subject != "principal-subject" {
				t.Fatalf("Verify() = %+v", verified)
			}
		})
	}
}

func TestJWTVerifierRejectsWrongAudienceExpiryAndUnknownKind(t *testing.T) {
	fixture := newJWTFixture(t)
	tests := map[string]string{
		"wrong audience": fixture.sign(t, string(PrincipalKindHuman), func(builder *jwt.Builder) {
			builder.Audience([]string{"another-api"})
		}),
		"expired": fixture.sign(t, string(PrincipalKindHuman), func(builder *jwt.Builder) {
			builder.Expiration(time.Now().Add(-time.Minute))
		}),
		"unknown kind": fixture.sign(t, "PLATFORM", nil),
	}
	for name, serialized := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := fixture.verifier.Verify(context.Background(), serialized); err == nil {
				t.Fatal("Verify() unexpectedly succeeded")
			}
		})
	}
}

func TestJWTVerifierRejectsTampering(t *testing.T) {
	fixture := newJWTFixture(t)
	serialized := fixture.sign(t, string(PrincipalKindHuman), nil)
	signatureStart := strings.LastIndex(serialized, ".") + 1
	replacement := byte('A')
	if serialized[signatureStart] == replacement {
		replacement = 'B'
	}
	serialized = serialized[:signatureStart] + string(replacement) + serialized[signatureStart+1:]
	if _, err := fixture.verifier.Verify(context.Background(), serialized); err == nil {
		t.Fatal("Verify() unexpectedly accepted a tampered token")
	}
}

func TestJWTVerifierUsesCurrentVerificationKeys(t *testing.T) {
	fixture := newJWTFixture(t)
	serialized := fixture.sign(t, string(PrincipalKindHuman), nil)
	if _, err := fixture.verifier.Verify(context.Background(), serialized); err != nil {
		t.Fatalf("Verify() before revocation error = %v", err)
	}
	if err := fixture.keys.RemoveKey(fixture.publicKey); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.verifier.Verify(context.Background(), serialized); err == nil {
		t.Fatal("Verify() accepted a token after its verification key was removed")
	}
}

func TestParseBearerAuthorization(t *testing.T) {
	token, err := ParseBearerAuthorization("Bearer signed-token")
	if err != nil || token != "signed-token" {
		t.Fatalf("ParseBearerAuthorization() = %q, %v", token, err)
	}
	for _, header := range []string{"", "Basic abc", "Bearer", "Bearer a b"} {
		if _, err := ParseBearerAuthorization(header); err == nil {
			t.Fatalf("ParseBearerAuthorization(%q) unexpectedly succeeded", header)
		}
	}
}
