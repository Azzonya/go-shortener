// Package pg provides a PostgreSQL implementation of the repository interface
// for managing shortened URLs.
package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

// St represents the PostgreSQL storage structure for shortened URLs.
type St struct {
	db *pgxpool.Pool // PostgreSQL connection pool.
}

// New creates and initializes a new PostgreSQL storage instance with the provided database connection pool.
func New(db *pgxpool.Pool) *St {
	s := &St{
		db: db,
	}

	var err error
	ctx := context.Background()

	if !s.TableExist(ctx) {
		err = s.Initialize(ctx)
		if err != nil {
			return nil
		}
	}

	return &St{
		db: db,
	}
}

// Initialize initializes the PostgreSQL storage by creating the necessary table if it doesn't exist.
func (s *St) Initialize(ctx context.Context) error {
	query := `CREATE TABLE urls (
				id SERIAL PRIMARY KEY,
				originalURL VARCHAR(255) NOT NULL,
				shortURL VARCHAR(255) UNIQUE NOT NULL,
                userID VARCHAR(255) NOT NULL,
                deleted BOOLEAN default false 
				);
				DO $$ 
				BEGIN 
   	 			IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE tablename = 'urls' AND indexname = 'idx_original_url') THEN
        			CREATE UNIQUE INDEX idx_original_url ON urls(originalURL);
    			END IF;
				END $$;`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("URL table creating error: %w", err)
	}

	return nil
}

// TableExist checks if the table exists in the PostgreSQL storage.
func (s *St) TableExist(ctx context.Context) bool {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) from urls").Scan(&count)

	return err == nil
}

// Add adds a new URL mapping to the PostgreSQL storage.
func (s *St) Add(ctx context.Context, originalURL, shortURL, userID string) error {
	query := `INSERT INTO urls (originalURL, shortURL, userID) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(ctx, query, originalURL, shortURL, userID)

	if err != nil {
		return err
	}

	return nil
}

// Update updates the short URL associated with the given original URL in the PostgreSQL storage.
func (s *St) Update(ctx context.Context, originalURL, shortURL string) error {
	query := `UPDATE urls SET shortURL = $1 where originalURL = $2`

	_, err := s.db.Exec(ctx, query, shortURL, originalURL)

	return err
}

// GetByShortURL retrieves the original URL associated with the given short URL from the PostgreSQL storage.
func (s *St) GetByShortURL(ctx context.Context, shortURL string) (string, bool) {
	var url string

	exist := false

	query := `SELECT originalURL from urls WHERE shortURL = $1`

	row := s.db.QueryRow(ctx, query, shortURL)

	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	if url != "" {
		exist = true
	}

	return url, exist
}

// GetByOriginalURL retrieves the short URL associated with the given original URL from the PostgreSQL storage.
func (s *St) GetByOriginalURL(ctx context.Context, originalURL string) (string, bool) {
	var url string

	exist := false

	query := `SELECT shortURL from urls WHERE originalURL = $1`

	row := s.db.QueryRow(ctx, query, originalURL)

	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	if url != "" {
		exist = true
	}

	return url, exist
}

// ListAll retrieves all shortened URLs associated with a user from the PostgreSQL storage.
func (s *St) ListAll(ctx context.Context, userID string) ([]*entities.ReqListAll, error) {
	result := []*entities.ReqListAll{}
	query := `SELECT originalurl, shortURL from urls WHERE userID = $1`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		URL := entities.ReqListAll{}
		err = rows.Scan(&URL.OriginalURL, &URL.ShortURL)
		if err != nil {
			return nil, err
		}

		result = append(result, &URL)
	}

	return result, nil
}

// CreateShortURLs creates multiple shortened URLs in a single transaction in the PostgreSQL storage.
func (s *St) CreateShortURLs(ctx context.Context, urls []*entities.ReqURL, userID string) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("error rollback: %d", err)
		}
	}()

	stmt, err := tx.Prepare(ctx, "insertURLs", "INSERT INTO urls (originalURL, shortURL, userID) VALUES($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("tx query error: %w", err)
	}

	for _, v := range urls {
		_, err := tx.Exec(ctx, stmt.Name, v.OriginalURL, v.ShortURL, userID)
		if err != nil {
			return fmt.Errorf("statement exec error: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

// DeleteURLs deletes URLs associated with the given userID.
// It takes a slice of URLs to delete and the userID of the user.
// It returns an error if any occurs during the deletion process.
func (s *St) DeleteURLs(ctx context.Context, urls []string, userID string) error {
	batch := &pgx.Batch{}

	for _, data := range urls {
		batch.Queue("UPDATE urls SET deleted = true WHERE shorturl = $1 AND userid = $2", data, userID)
	}

	bRes := s.db.SendBatch(ctx, batch)
	err := bRes.Close()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	return err
}

// URLDeleted checks if the URL with the given shortURL is deleted.
func (s *St) URLDeleted(ctx context.Context, shortURL string) bool {
	deleted := false
	query := `SELECT deleted from urls WHERE shortURL = $1`

	row := s.db.QueryRow(ctx, query, shortURL)

	err := row.Scan(&deleted)
	if err != nil {
		return false
	}

	return deleted
}

// WriteEvent writes an event to the storage.
func (s *St) WriteEvent(_ *inmemory.Event) error {
	return nil
}

// SyncData synchronizes data.
// This method does not take any input parameters or return anything.
func (s *St) SyncData() {
	//
}

// Ping pings the database.
func (s *St) Ping(ctx context.Context) error {
	err := s.db.Ping(ctx)

	return err
}

// CountUsers returns the number of unique users (userid) from the urls table.
// It returns the count of unique userids and an error, if one occurred during the query execution.
func (s *St) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(DISTINCT userid) FROM urls").Scan(&count)

	return count, err
}

// CountURLs returns the number of unique URLs (originalurl) from the urls table.
// It returns the count of unique originalurls and an error, if one occurred during the query execution.
func (s *St) CountURLs(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(DISTINCT originalurl) FROM urls").Scan(&count)

	return count, err
}
