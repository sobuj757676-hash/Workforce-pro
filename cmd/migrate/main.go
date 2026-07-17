package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/workforce-pro/workforce-payroll/internal/persistence"
)

func main() {
	if err := run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(arguments []string) error {
	if len(arguments) != 2 || (arguments[1] != "up" && arguments[1] != "down") {
		return errors.New("usage: migrate up|down")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	migrator, err := persistence.NewMigrator(pool)
	if err != nil {
		return err
	}
	if arguments[1] == "up" {
		return migrator.Up(ctx)
	}
	return migrator.Down(ctx)
}
