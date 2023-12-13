package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type St struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *St {
	s := &St{
		db: db,
	}

	var err error
	exist := s.URLTableExist()

	if !exist {
		err = s.URLTableInit()
		if err != nil {
			return nil
		}
	}

	return &St{
		db: db,
	}
}

func (s *St) URLTableInit() error {
	query := `CREATE TABLE urls (
				id SERIAL PRIMARY KEY,
				originalURL VARCHAR(255) NOT NULL,
				shortURL VARCHAR(255) UNIQUE NOT NULL
				);`

	_, err := s.db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("URL table creating error: %w", err)
	}

	return nil
}

func (s *St) URLTableExist() bool {
	var count int
	err := s.db.QueryRow(context.Background(), "SELECT COUNT(*) from urls").Scan(&count)
	if err == nil {
		return true
	}

	return false
}

func (s *St) URLAddNew(originalURL, shortURL string) error {
	query := `INSERT INTO urls (originalURL, shortURL) VALUES ($1, $2)`

	_, err := s.db.Exec(context.Background(), query, originalURL, shortURL)
	if err != nil {
		return fmt.Errorf("error: db add new url line: %w", err)
	}

	return nil
}

func (s *St) URLGet(shortURL string) (string, bool) {
	var url string

	exist := false

	query := `SELECT originalURL from urls WHERE shortURL = $1`

	row := s.db.QueryRow(context.Background(), query, shortURL)

	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	if url != "" {
		exist = true
	}

	return url, exist
}

func (s *St) Ping() error {
	err := s.db.Ping(context.Background())

	return err
}
