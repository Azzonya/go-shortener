package repo

import (
	_ "github.com/jackc/pgx/v5"
)

type Repo interface {
	URLTableInit() error
	URLTableExist() error
	URLAddNew(originalURL, shortURL string) error
	URLGet(shortURL string) (string, bool)

	Ping() error
}
