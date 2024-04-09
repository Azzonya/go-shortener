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

type St struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *St {
	s := &St{
		db: db,
	}

	var err error

	if !s.TableExist() {
		err = s.Initialize()
		if err != nil {
			return nil
		}
	}

	return &St{
		db: db,
	}
}

func (s *St) Initialize() error {
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

	_, err := s.db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("URL table creating error: %w", err)
	}

	return nil
}

func (s *St) TableExist() bool {
	var count int
	err := s.db.QueryRow(context.Background(), "SELECT COUNT(*) from urls").Scan(&count)

	return err == nil
}

func (s *St) Add(originalURL, shortURL, userID string) error {
	query := `INSERT INTO urls (originalURL, shortURL, userID) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(context.Background(), query, originalURL, shortURL, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *St) Update(originalURL, shortURL string) error {
	query := `UPDATE urls SET shortURL = $1 where originalURL = $2`

	_, err := s.db.Exec(context.Background(), query, shortURL, originalURL)

	return err
}

func (s *St) GetByShortURL(shortURL string) (string, bool) {
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

func (s *St) GetByOriginalURL(originalURL string) (string, bool) {
	var url string

	exist := false

	query := `SELECT shortURL from urls WHERE originalURL = $1`

	row := s.db.QueryRow(context.Background(), query, originalURL)

	err := row.Scan(&url)
	if err != nil {
		return "", false
	}

	if url != "" {
		exist = true
	}

	return url, exist
}

func (s *St) ListAll(userID string) ([]*entities.ReqListAll, error) {
	result := []*entities.ReqListAll{}
	query := `SELECT originalurl, shortURL from urls WHERE userID = $1`

	rows, err := s.db.Query(context.Background(), query, userID)
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

func (s *St) CreateShortURLs(urls []*entities.ReqURL, userID string) error {
	ctx := context.Background()

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

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	return nil
}

func (s *St) DeleteURLs(urls []string, userID string) error {
	batch := &pgx.Batch{}

	for _, data := range urls {
		batch.Queue("UPDATE urls SET deleted = true WHERE shorturl = $1 AND userid = $2", data, userID)
	}

	bRes := s.db.SendBatch(context.Background(), batch)
	err := bRes.Close()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	return err
}

func (s *St) URLDeleted(shortURL string) bool {
	deleted := false
	query := `SELECT deleted from urls WHERE shortURL = $1`

	row := s.db.QueryRow(context.Background(), query, shortURL)

	err := row.Scan(&deleted)
	if err != nil {
		return false
	}

	return deleted
}

func (s *St) WriteEvent(event *inmemory.Event) error {
	return nil
}

func (s *St) SyncData() {
	//
}

func (s *St) Ping() error {
	err := s.db.Ping(context.Background())

	return err
}
