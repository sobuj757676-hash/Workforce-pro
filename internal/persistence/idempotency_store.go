package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/workforce-pro/workforce-payroll/internal/foundation"
)

// PostgresIdempotencyStore implements foundation.IdempotencyStore using
// the idempotency_records table.
// Requirements: 7.1–7.4
type PostgresIdempotencyStore struct {
	pool *pgxpool.Pool
}

// NewPostgresIdempotencyStore creates a new idempotency store.
func NewPostgresIdempotencyStore(pool *pgxpool.Pool) (*PostgresIdempotencyStore, error) {
	if pool == nil {
		return nil, errors.New("postgres pool is required")
	}
	return &PostgresIdempotencyStore{pool: pool}, nil
}

// Find returns an existing idempotency record for the given Tenant and key.
func (s *PostgresIdempotencyStore) Find(ctx context.Context, tenantID foundation.Identifier, key string) (foundation.IdempotencyRecord, error) {
	var record foundation.IdempotencyRecord
	var digestText, resultRef string
	var createdAt, expiresAt time.Time

	err := s.pool.QueryRow(ctx, `
		SELECT request_digest, result_ref, created_at, expires_at
		FROM idempotency_records
		WHERE tenant_id = $1 AND idempotency_key = $2 AND expires_at > CURRENT_TIMESTAMP
	`, tenantID.String(), key).Scan(&digestText, &resultRef, &createdAt, &expiresAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return foundation.IdempotencyRecord{}, foundation.ErrNotFound
	}
	if err != nil {
		return foundation.IdempotencyRecord{}, err
	}

	digest, err := foundation.ParseDigest(digestText)
	if err != nil {
		return foundation.IdempotencyRecord{}, err
	}

	record.TenantID = tenantID
	record.Key = key
	record.RequestDigest = digest
	record.ResultRef = resultRef
	record.CreatedAt = createdAt
	record.ExpiresAt = expiresAt
	return record, nil
}

// Store persists a new idempotency record. If a record with the same
// Tenant+key already exists, it returns ErrConflict.
func (s *PostgresIdempotencyStore) Store(ctx context.Context, record foundation.IdempotencyRecord) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO idempotency_records (tenant_id, idempotency_key, request_digest, result_ref, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (tenant_id, idempotency_key) DO NOTHING
	`, record.TenantID.String(), record.Key, record.RequestDigest.String(), record.ResultRef, record.CreatedAt, record.ExpiresAt)
	if err != nil {
		return err
	}
	// Check if the insert actually happened by reading back.
	existing, err := s.Find(ctx, record.TenantID, record.Key)
	if err != nil {
		return err
	}
	if !existing.RequestDigest.Equal(record.RequestDigest) {
		return foundation.ErrIdempotencyKeyReuse
	}
	return nil
}
