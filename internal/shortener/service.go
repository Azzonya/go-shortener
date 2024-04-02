package shortener

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/repo"
	"github.com/google/uuid"
	"net/url"
)

type Shortener struct {
	repo    repo.Repo
	baseURL string
}

func New(baseURL string, repo repo.Repo) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		repo:    repo,
	}
}

func (s *Shortener) GetOneByShortURL(key string) (string, bool) {
	return s.repo.GetByShortURL(key)
}

func (s *Shortener) GetOneByOriginalURL(url string) (string, bool) {
	URL, exist := s.repo.GetByOriginalURL(url)
	if !exist {
		return "", false
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, URL)

	return outputURL, exist
}

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

func (s *Shortener) ShortenAndSaveLink(originalURL, userID string) (string, error) {
	shortURL := s.GenerateShortURL()
	if err := s.repo.Add(originalURL, shortURL, userID); err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

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

func (s *Shortener) DeleteURLs(urls []string, userID string) {
	go func() {
		if err := s.repo.DeleteURLs(urls, userID); err != nil {
			logger.Log.Error("Failed to delete URLs " + err.Error())
		}
	}()
}

func (s *Shortener) IsDeleted(shortURL string) bool {
	return s.repo.URLDeleted(shortURL)
}

func (s *Shortener) GenerateShortURL() string {
	return uuid.New().String()
}

func (s *Shortener) PingDB() error {
	err := s.repo.Ping()

	return err
}
