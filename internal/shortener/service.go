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
	UseDB   bool
}

func New(baseURL string, repo repo.Repo, UseDB bool) *Shortener {
	return &Shortener{
		baseURL: baseURL,
		repo:    repo,
		UseDB:   UseDB,
	}
}

func (s *Shortener) GetOneByShortURL(key string) (string, bool) {
	return s.repo.GetByShortURL(key)
}

func (s *Shortener) GetOneByOriginalURL(url string) (string, bool) {
	URL, exist := s.repo.GetByOriginalURL(url)

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
	var err error
	shortURL := s.GenerateShortURL()

	err = s.repo.Add(originalURL, shortURL, userID)

	if err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s", s.baseURL, shortURL)

	return outputURL, nil
}

func (s *Shortener) ShortenURLs(urls []*entities.ReqURL, userID string) ([]*entities.ReqURL, error) {
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
