package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ffaeso/arctic/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresPool creates a PostgreSQL connection pool configured for Arctic.
//
// The returned *sql.DB manages a pool of connections and is safe for concurrent use.
// The caller is responsible for calling Close() when the pool is no longer needed.
//
// Returns an error if the connection cannot be established or the database is unreachable.
func NewPostgresPool(cfg *config.DatasourceConfig) (*sql.DB, error) {
	conn, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}

	// pool configurations - we are keeping it generic right now
	// without configuration options unless requested
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)
	conn.SetConnMaxIdleTime(1 * time.Minute)

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	return conn, err
}
