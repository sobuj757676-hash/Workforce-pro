package config

import "testing"

func TestFromEnvironmentRequiresCompleteSecureConfiguration(t *testing.T) {
	t.Setenv("HTTP_ADDRESS", ":8080")
	t.Setenv("DATABASE_URL", "postgres://service:secret@database/workforce")
	t.Setenv("OIDC_ISSUER", "https://identity.example.test/realms/workforce")
	t.Setenv("OIDC_AUDIENCE", "workforce-api")
	t.Setenv("OIDC_JWKS_URL", "https://identity.example.test/realms/workforce/protocol/openid-connect/certs")
	t.Setenv("OIDC_ACCEPTABLE_SKEW", "30s")
	t.Setenv("OIDC_JWKS_REFRESH_INTERVAL", "5m")

	config, err := FromEnvironment()
	if err != nil {
		t.Fatalf("FromEnvironment() error = %v", err)
	}
	if config.OIDCAudience != "workforce-api" || config.OIDCAcceptableSkew.String() != "30s" {
		t.Fatalf("FromEnvironment() = %+v", config)
	}
}

func TestFromEnvironmentRejectsInsecureIdentityURLsAndMissingSkew(t *testing.T) {
	setValidEnvironment := func(t *testing.T) {
		t.Helper()
		t.Setenv("HTTP_ADDRESS", ":8080")
		t.Setenv("DATABASE_URL", "postgres://service:secret@database/workforce")
		t.Setenv("OIDC_ISSUER", "https://identity.example.test/realms/workforce")
		t.Setenv("OIDC_AUDIENCE", "workforce-api")
		t.Setenv("OIDC_JWKS_URL", "https://identity.example.test/certs")
		t.Setenv("OIDC_ACCEPTABLE_SKEW", "30s")
		t.Setenv("OIDC_JWKS_REFRESH_INTERVAL", "5m")
	}

	t.Run("insecure issuer", func(t *testing.T) {
		setValidEnvironment(t)
		t.Setenv("OIDC_ISSUER", "http://identity.example.test")
		if _, err := FromEnvironment(); err == nil {
			t.Fatal("FromEnvironment() unexpectedly accepted HTTP issuer")
		}
	})
	t.Run("missing skew", func(t *testing.T) {
		setValidEnvironment(t)
		t.Setenv("OIDC_ACCEPTABLE_SKEW", "")
		if _, err := FromEnvironment(); err == nil {
			t.Fatal("FromEnvironment() unexpectedly accepted missing skew")
		}
	})
}
