package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/workforce-pro/workforce-payroll/internal/foundation"
)

var ErrPrincipalNotFound = errors.New("active principal not found")

type PostgresPrincipalRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPrincipalRepository(pool *pgxpool.Pool) (*PostgresPrincipalRepository, error) {
	if pool == nil {
		return nil, errors.New("postgres pool is required")
	}
	return &PostgresPrincipalRepository{pool: pool}, nil
}

func (r *PostgresPrincipalRepository) ResolveActive(ctx context.Context, verified VerifiedIdentity) (Principal, error) {
	if err := verified.Validate(); err != nil {
		return Principal{}, ErrPrincipalNotFound
	}

	var principalIDText, tenantIDText, issuer, subject, kind string
	err := r.pool.QueryRow(ctx, `
		SELECT principal_id::text, tenant_id::text, issuer, subject, principal_kind
		FROM principal_bindings
		WHERE issuer = $1
		  AND subject = $2
		  AND principal_kind = $3
		  AND status = 'ACTIVE'
	`, verified.Issuer, verified.Subject, string(verified.Kind)).Scan(
		&principalIDText, &tenantIDText, &issuer, &subject, &kind,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrPrincipalNotFound
	}
	if err != nil {
		return Principal{}, fmt.Errorf("resolve active principal: %w", err)
	}

	principalID, err := foundation.ParseIdentifier(principalIDText)
	if err != nil {
		return Principal{}, fmt.Errorf("invalid persisted principal identifier: %w", err)
	}
	tenantID, err := foundation.ParseIdentifier(tenantIDText)
	if err != nil {
		return Principal{}, fmt.Errorf("invalid persisted tenant identifier: %w", err)
	}

	principal := Principal{
		ID:       principalID,
		TenantID: tenantID,
		Issuer:   issuer,
		Subject:  subject,
		Kind:     PrincipalKind(kind),
	}
	if err := principal.Validate(); err != nil {
		return Principal{}, fmt.Errorf("invalid persisted principal: %w", err)
	}
	return principal, nil
}
