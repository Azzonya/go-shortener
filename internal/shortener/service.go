package shortener

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/repo"
	"github.com/Azzonya/go-shortener/internal/storage"
	"math/rand"
	"time"
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
	const shorURLLenth = 8

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, shorURLLenth)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Shortener) PingDB() error {
	err := s.repo.Ping()

	return err
}
