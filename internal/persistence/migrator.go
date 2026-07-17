package persistence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/workforce-pro/workforce-payroll/migrations"
)

const principalBindingsVersion int64 = 1

type Migrator struct {
	pool *pgxpool.Pool
}

func NewMigrator(pool *pgxpool.Pool) (*Migrator, error) {
	if pool == nil {
		return nil, errors.New("postgres pool is required")
	}
	return &Migrator{pool: pool}, nil
}

func (m *Migrator) Up(ctx context.Context) error {
	checksum := migrationChecksum(migrations.PrincipalBindingsUp)
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", principalBindingsVersion); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	var existingChecksum string
	err = tx.QueryRow(ctx, "SELECT checksum FROM schema_migrations WHERE version = $1", principalBindingsVersion).Scan(&existingChecksum)
	if err == nil {
		if existingChecksum != checksum {
			return errors.New("applied migration checksum mismatch")
		}
		return tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if _, err := tx.Exec(ctx, migrations.PrincipalBindingsUp); err != nil {
		return fmt.Errorf("apply migration 1: %w", err)
	}
	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, checksum) VALUES ($1, $2)", principalBindingsVersion, checksum); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *Migrator) Down(ctx context.Context) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", principalBindingsVersion); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, migrations.PrincipalBindingsDown); err != nil {
		return fmt.Errorf("revert migration 1: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", principalBindingsVersion); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func migrationChecksum(sql string) string {
	digest := sha256.Sum256([]byte(sql))
	return hex.EncodeToString(digest[:])
}
