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

type migrationEntry struct {
	version int64
	sql     string
	downSQL string
}

var allMigrations = []migrationEntry{
	{version: 1, sql: migrations.PrincipalBindingsUp, downSQL: migrations.PrincipalBindingsDown},
	{version: 2, sql: migrations.IdempotencyRecordsUp, downSQL: migrations.IdempotencyRecordsDown},
}

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
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", int64(9999)); err != nil {
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

	for _, migration := range allMigrations {
		checksum := migrationChecksum(migration.sql)
		var existingChecksum string
		err := tx.QueryRow(ctx, "SELECT checksum FROM schema_migrations WHERE version = $1", migration.version).Scan(&existingChecksum)
		if err == nil {
			if existingChecksum != checksum {
				return fmt.Errorf("applied migration %d checksum mismatch", migration.version)
			}
			continue
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if _, err := tx.Exec(ctx, migration.sql); err != nil {
			return fmt.Errorf("apply migration %d: %w", migration.version, err)
		}
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, checksum) VALUES ($1, $2)", migration.version, checksum); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (m *Migrator) Down(ctx context.Context) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", int64(9999)); err != nil {
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

	// Apply down migrations in reverse order.
	for i := len(allMigrations) - 1; i >= 0; i-- {
		migration := allMigrations[i]
		if _, err := tx.Exec(ctx, migration.downSQL); err != nil {
			return fmt.Errorf("revert migration %d: %w", migration.version, err)
		}
		if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", migration.version); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func migrationChecksum(sql string) string {
	digest := sha256.Sum256([]byte(sql))
	return hex.EncodeToString(digest[:])
}
