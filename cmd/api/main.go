package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/workforce-pro/workforce-payroll/internal/config"
	"github.com/workforce-pro/workforce-payroll/internal/httpapi"
	"github.com/workforce-pro/workforce-payroll/internal/identity"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.FromEnvironment()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	keyCache, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		return err
	}
	defer func() { _ = keyCache.Shutdown(context.Background()) }()
	if err := keyCache.Register(ctx, cfg.OIDCJWKSURL,
		jwk.WithMinInterval(cfg.OIDCJWKSRefresh),
		jwk.WithMaxInterval(cfg.OIDCJWKSRefresh),
	); err != nil {
		return err
	}
	keys, err := keyCache.CachedSet(cfg.OIDCJWKSURL)
	if err != nil {
		return err
	}
	if _, err := keyCache.Refresh(ctx, cfg.OIDCJWKSURL); err != nil {
		return err
	}
	verifier, err := identity.NewJWTVerifier(cfg.OIDCIssuer, cfg.OIDCAudience, keys, cfg.OIDCAcceptableSkew)
	if err != nil {
		return err
	}
	principals, err := identity.NewPostgresPrincipalRepository(pool)
	if err != nil {
		return err
	}
	authenticator, err := identity.NewAuthenticator(verifier, principals)
	if err != nil {
		return err
	}
	handler, err := httpapi.NewServer(authenticator)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errChannel := make(chan error, 1)
	go func() {
		errChannel <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownContext)
	case err := <-errChannel:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
