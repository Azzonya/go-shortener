package shortener

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/repo"
	"github.com/Azzonya/go-shortener/internal/storage"
	"github.com/google/uuid"
)

type Shortener struct {
	repo    repo.Repo
	storage *storage.Storage
	baseURL string
	UseDB   bool
}

func New(baseURL string, storage *storage.Storage, repo repo.Repo, UseDB bool) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		storage: storage,
		repo:    repo,
		UseDB:   UseDB,
	}
}

func (s *Shortener) GetOne(key string) (string, bool) {
	var URL string
	var exist bool

	if !s.UseDB {
		URL, exist = s.storage.GetOne(key)
	} else {
		URL, exist = s.repo.URLGetByShortURL(key)
	}

	return URL, exist
}

func (s *Shortener) ShortenAndSaveLink(originalURL string) (string, error) {
	var err error
	shortURL := s.GenerateShortURL()

	if !s.UseDB {
		err = s.storage.Add(shortURL, originalURL)
	} else {
		err = s.repo.URLAddNew(originalURL, shortURL)
	}

	if err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

func (s *Shortener) GenerateShortURL() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *Shortener) PingDB() error {
	err := s.repo.Ping()

	return err
}
