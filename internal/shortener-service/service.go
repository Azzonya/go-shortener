package shortener_service

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/storage"
	"math/rand"
)

type Shortener struct {
	baseURL string
	storage *storage.Storage
}

func New(baseURL string, storage *storage.Storage) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		storage: storage,
	}
}

func (s *Shortener) GetOne(key string) (string, bool) {
	URL, exist := s.storage.GetOne(key)
	return URL, exist
}

func (s *Shortener) ShortenAndSaveLink(originalURL string) (string, error) {
	shortURL := s.GenerateShortURL()

	err := s.storage.Add(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

func (s *Shortener) GenerateShortURL() string {
	const shorURLLenth = 8
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, shorURLLenth)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
