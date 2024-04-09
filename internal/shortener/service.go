// Package shortener provides functionalities for URL shortening, storage, retrieval, and management.
// It includes methods for shortening long URLs, retrieving original URLs from short URLs,
// listing shortened URLs associated with a user ID, deleting specific shortened URLs,
// checking the status of a short URL, interacting with the database for storage and retrieval,
// and handling concurrency for URL deletion.
//
// The package is designed to offer a comprehensive solution for URL shortening needs,
// ensuring robustness, performance, and usability.
package shortener

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/repo"
)

// Shortener represents the URL shortener service.
type Shortener struct {
	repo    repo.Repo // Repository for storing and retrieving shortened URLs
	baseURL string    // Base URL used for constructing shortened URLs
}

// New creates a new instance of the Shortener struct.
func New(baseURL string, repo repo.Repo) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		repo:    repo,
	}
}

// GetOneByShortURL retrieves the original URL associated with a given short URL.
func (s *Shortener) GetOneByShortURL(key string) (string, bool) {
	return s.repo.GetByShortURL(key)
}

// GetOneByOriginalURL retrieves the short URL associated with a given original URL.
func (s *Shortener) GetOneByOriginalURL(url string) (string, bool) {
	URL, exist := s.repo.GetByOriginalURL(url)
	if !exist {
		return "", false
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, URL)

	return outputURL, exist
}

// ListAll retrieves a list of all shortened URLs associated with a given user ID.
func (s *Shortener) ListAll(userID string) ([]*entities.ReqListAll, error) {
	list, err := s.repo.ListAll(userID)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		v.ShortURL = fmt.Sprintf("%s/%s", s.baseURL, v.ShortURL)
	}
	return list, nil
}

// ShortenAndSaveLink generates a short URL for a given original URL and saves it in the repository.
func (s *Shortener) ShortenAndSaveLink(originalURL, userID string) (string, error) {
	shortURL := s.GenerateShortURL()
	if err := s.repo.Add(originalURL, shortURL, userID); err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

// ShortenURLs shortens multiple URLs simultaneously and saves them in the repository.
func (s *Shortener) ShortenURLs(urls []*entities.ReqURL, userID string) ([]*entities.ReqURL, error) {
	shortenedURLs := make([]*entities.ReqURL, len(urls))

	for i, u := range urls {
		shortURL := s.GenerateShortURL()

		shortenedURLs[i] = &entities.ReqURL{
			ID:          u.ID,
			OriginalURL: u.OriginalURL,
			ShortURL:    shortURL,
		}

		urls[i].ShortURL = shortURL
	}

	err := s.repo.CreateShortURLs(shortenedURLs, userID)
	if err != nil {
		return nil, err
	}

	for i, v := range shortenedURLs {
		resultString, err := url.JoinPath(s.baseURL, v.ShortURL)
		if err != nil {
			return nil, err
		}

		shortenedURLs[i].ShortURL = resultString
		shortenedURLs[i].OriginalURL = ""
	}

	return shortenedURLs, nil
}

// DeleteURLs deletes URLs associated with a given user ID.
func (s *Shortener) DeleteURLs(urls []string, userID string) {
	go func() {
		if err := s.repo.DeleteURLs(urls, userID); err != nil {
			logger.Log.Error("Failed to delete URLs " + err.Error())
		}
	}()
}

// IsDeleted checks if a given short URL has been deleted.
func (s *Shortener) IsDeleted(shortURL string) bool {
	return s.repo.URLDeleted(shortURL)
}

// GenerateShortURL generates a unique short URL using UUID.
func (s *Shortener) GenerateShortURL() string {
	return uuid.New().String()
}

// PingDB pings the database to check its connectivity.
func (s *Shortener) PingDB() error {
	err := s.repo.Ping()

	return err
}
