package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type St struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *St {
	s := &St{
		db: db,
	}

	var err error
	exist := s.TableExist()

	if !exist {
		err = s.TableInit()
		if err != nil {
			return nil
		}
	}

	return &St{
		db: db,
	}
}

func (s *St) TableInit() error {
	query := `CREATE TABLE urls (
				id SERIAL PRIMARY KEY,
				originalURL VARCHAR(255) NOT NULL,
				shortURL VARCHAR(255) UNIQUE NOT NULL
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

func (s *St) AddNew(originalURL, shortURL string) error {
	query := `INSERT INTO urls (originalURL, shortURL) VALUES ($1, $2)`

	err := s.db.QueryRow(context.Background(), query, originalURL, shortURL)

	if err != nil {
		return errors.New("cannot insert line")
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

func (s *St) CreateShortURLs(urls []*entities.ReqURL) error {
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

	stmt, err := tx.Prepare(ctx, "insertURLs", "INSERT INTO urls (originalURL, shortURL) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("tx query error: %w", err)
	}

	for _, v := range urls {
		_, err := tx.Exec(ctx, stmt.Name, v.OriginalURL, v.ShortURL)
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

func (s *St) Ping() error {
	err := s.db.Ping(context.Background())

	return err
}
