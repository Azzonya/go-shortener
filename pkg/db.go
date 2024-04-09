// Package pkg provides utility functions for working with databases.
//
// This package includes functions for initializing database connections and executing queries.
package pkg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDatabasePg initializes a connection pool to a PostgreSQL database.
// It takes a PostgreSQL connection string (pgDsn) as input and returns a *pgxpool.Pool representing the database connection pool.
// If initialization fails, it returns an error.
func InitDatabasePg(pgDsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(pgDsn)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
