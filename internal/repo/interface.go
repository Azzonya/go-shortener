package repo

import (
	_ "github.com/jackc/pgx/v5"
)

type Repo interface {
	URLTableInit() error
	URLTableExist() bool
	URLAddNew(originalURL, shortURL string) error
	URLUpdate(originalURL, shortURL string) error
	URLGetByShortURL(shortURL string) (string, bool)
	URLGetByOriginalURL(originalURL string) (string, bool)

	Ping() error
}
