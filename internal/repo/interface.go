package repo

import (
	"github.com/Azzonya/go-shortener/internal/entities"
	_ "github.com/jackc/pgx/v5"
)

type Repo interface {
	TableInit() error
	TableExist() bool
	AddNew(originalURL, shortURL, userID string) error
	CreateShortURLs(urls []*entities.ReqURL, userID string) error
	Update(originalURL, shortURL string) error
	GetByShortURL(shortURL string) (string, bool)
	GetByOriginalURL(originalURL string) (string, bool)
	ListAll(userID string) ([]*entities.ReqListAll, error)
	DeleteURLs(urls []string, userID string) error
	URLDeleted(shortURL string) bool

	Ping() error
}
