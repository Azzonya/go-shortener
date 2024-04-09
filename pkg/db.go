package pkg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
