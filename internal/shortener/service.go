package shortener

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/inmemory"
	"github.com/Azzonya/go-shortener/internal/repo"
	"github.com/google/uuid"
	"net/url"
)

type Shortener struct {
	repo     repo.Repo
	inmemory *inmemory.Storage
	baseURL  string
	UseDB    bool
	UserID   string
}

func New(baseURL string, inmemory *inmemory.Storage, repo repo.Repo, UserID string, UseDB bool) *Shortener {
	return &Shortener{
		baseURL:  baseURL,
		inmemory: inmemory,
		repo:     repo,
		UserID:   UserID,
		UseDB:    UseDB,
	}
}

func (s *Shortener) GetOneByShortURL(key string) (string, bool) {
	var URL string
	var exist bool

	if s.UseDB {
		URL, exist = s.repo.GetByShortURL(key)
	} else {
		URL, exist = s.inmemory.GetOne(key)
	}

	return URL, exist
}

func (s *Shortener) GetOneByOriginalURL(url string) (string, bool) {
	URL, exist := s.repo.GetByOriginalURL(url)

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, URL)

	return outputURL, exist
}

func (s *Shortener) ListAll() ([]*entities.ReqListAll, error) {
	list, err := s.repo.ListAll(s.UserID)
	if err != nil {
		return nil, err
	}

	for _, v := range list {
		v.ShortURL = fmt.Sprintf("%s/%s", s.baseURL, v.ShortURL)
	}
	return list, nil
}

func (s *Shortener) ShortenAndSaveLink(originalURL string) (string, error) {
	var err error
	shortURL := s.GenerateShortURL()

	if !s.UseDB {
		err = s.inmemory.Add(shortURL, originalURL)
	} else {
		err = s.repo.AddNew(originalURL, shortURL, s.UserID)
	}

	if err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

func (s *Shortener) ShortenURLs(urls []*entities.ReqURL) ([]*entities.ReqURL, error) {
	var shortenedURLs []*entities.ReqURL

	for i := range urls {
		shortURL := s.GenerateShortURL()
		urls[i].ShortURL = shortURL

		shortenedURLs = append(shortenedURLs, &entities.ReqURL{
			ID:          urls[i].ID,
			OriginalURL: urls[i].OriginalURL,
			ShortURL:    shortURL,
		})
	}

	err := s.repo.CreateShortURLs(shortenedURLs, s.UserID)
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

func (s *Shortener) DeleteURLs(urls []string) error {
	return s.repo.DeleteURLs(urls, s.UserID)
}

func (s *Shortener) IsDeleted(shortURL string) bool {
	if !s.UseDB {
		return false
	}
	return s.repo.URLDeleted(shortURL)
}

func (s *Shortener) GenerateShortURL() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func (s *Shortener) PingDB() error {
	err := s.repo.Ping()

	return err
}
