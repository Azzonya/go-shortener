package repo

import (
	_ "github.com/jackc/pgx/v5"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

type Repo interface {
	Initialize() error
	TableExist() bool
	Add(originalURL, shortURL, userID string) error
	CreateShortURLs(urls []*entities.ReqURL, userID string) error
	Update(originalURL, shortURL string) error
	GetByShortURL(shortURL string) (string, bool)
	GetByOriginalURL(originalURL string) (string, bool)
	ListAll(userID string) ([]*entities.ReqListAll, error)
	DeleteURLs(urls []string, userID string) error
	URLDeleted(shortURL string) bool
	WriteEvent(event *inmemory.Event) error
	SyncData()

	Ping() error
}
