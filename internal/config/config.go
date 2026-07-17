package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPAddress        string
	DatabaseURL        string
	OIDCIssuer         string
	OIDCAudience       string
	OIDCJWKSURL        string
	OIDCAcceptableSkew time.Duration
	OIDCJWKSRefresh    time.Duration
}

func FromEnvironment() (Config, error) {
	config := Config{
		HTTPAddress:  strings.TrimSpace(os.Getenv("HTTP_ADDRESS")),
		DatabaseURL:  strings.TrimSpace(os.Getenv("DATABASE_URL")),
		OIDCIssuer:   strings.TrimSpace(os.Getenv("OIDC_ISSUER")),
		OIDCAudience: strings.TrimSpace(os.Getenv("OIDC_AUDIENCE")),
		OIDCJWKSURL:  strings.TrimSpace(os.Getenv("OIDC_JWKS_URL")),
	}
	if config.HTTPAddress == "" || config.DatabaseURL == "" || config.OIDCIssuer == "" || config.OIDCAudience == "" || config.OIDCJWKSURL == "" {
		return Config{}, errors.New("HTTP_ADDRESS, DATABASE_URL, OIDC_ISSUER, OIDC_AUDIENCE, and OIDC_JWKS_URL are required")
	}

	skewText := strings.TrimSpace(os.Getenv("OIDC_ACCEPTABLE_SKEW"))
	if skewText == "" {
		return Config{}, errors.New("OIDC_ACCEPTABLE_SKEW is required")
	}
	skew, err := time.ParseDuration(skewText)
	if err != nil || skew < 0 {
		return Config{}, errors.New("OIDC_ACCEPTABLE_SKEW must be a non-negative duration")
	}
	config.OIDCAcceptableSkew = skew

	refreshText := strings.TrimSpace(os.Getenv("OIDC_JWKS_REFRESH_INTERVAL"))
	if refreshText == "" {
		return Config{}, errors.New("OIDC_JWKS_REFRESH_INTERVAL is required")
	}
	refresh, err := time.ParseDuration(refreshText)
	if err != nil || refresh <= 0 {
		return Config{}, errors.New("OIDC_JWKS_REFRESH_INTERVAL must be a positive duration")
	}
	config.OIDCJWKSRefresh = refresh

	if err := requireHTTPS("OIDC_ISSUER", config.OIDCIssuer); err != nil {
		return Config{}, err
	}
	if err := requireHTTPS("OIDC_JWKS_URL", config.OIDCJWKSURL); err != nil {
		return Config{}, err
	}
	return config, nil
}

func requireHTTPS(name, value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("%s must be an absolute HTTPS URL without user information", name)
	}
	return nil
}
