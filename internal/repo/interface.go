package repo

import (
	"github.com/Azzonya/go-shortener/internal/entities"
	_ "github.com/jackc/pgx/v5"
)

type Repo interface {
	TableInit() error
	TableExist() bool
	AddNew(originalURL, shortURL string) error
	Update(originalURL, shortURL string) error
	GetByShortURL(shortURL string) (string, bool)
	GetByOriginalURL(originalURL string) (string, bool)
	CreateShortURLs(urls []*entities.ReqURL) error

	Ping() error
}
