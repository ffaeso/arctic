package persistence

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunMigrations applies all pending database migrations.
// Returns nil if migrations succeed or if there are no migrations to apply.
func RunMigrations(pool *sql.DB) error {
	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(pool, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
