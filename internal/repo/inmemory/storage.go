// Package inmemory provides an in-memory implementation of the repository interface
// for managing shortened URLs.
package inmemory

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Azzonya/go-shortener/internal/entities"
)

// St represents the in-memory storage structure for shortened URLs.
type St struct {
	URLMap   map[string]string // Map to store shortened URLs as key-value pairs.
	filePath string            // File path to store the JSON data.
	lastID   int               // Last ID used for the storage.
}

// Event represents the event structure used for JSON encoding and decoding.
type Event struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// New creates and initializes a new in-memory storage instance with the provided file path.
func New(filePath string) (*St, error) {
	s := &St{}

	s.filePath = filePath
	err := s.Initialize()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Initialize initializes the in-memory storage by reading data from the provided file.
func (s *St) Initialize() error {
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	newDecoder := json.NewDecoder(file)

	s.URLMap = make(map[string]string)
	s.lastID = 0

	for {
		var event Event
		if err := newDecoder.Decode(&event); err != nil {
			if err != io.EOF {
				log.Println("error decode JSON:", err)
			}
			break
		}
		s.lastID++
		s.URLMap[event.OriginalURL] = event.ShortURL
	}

	return nil
}

// TableExist checks if the table exists in the in-memory storage (always returns true).
func (s *St) TableExist() bool {
	return true
}

// Add adds a new URL mapping to the in-memory storage.
func (s *St) Add(originalURL, shortURL, _ string) error {
	s.URLMap[shortURL] = originalURL
	s.lastID++

	return nil
}

// Update always returns nil.
func (s *St) Update(_, _ string) error {
	return nil
}

// GetByShortURL retrieves the original URL associated with the given short URL.
func (s *St) GetByShortURL(shortURL string) (string, bool) {
	URL, exist := s.URLMap[shortURL]
	return URL, exist
}

// GetByOriginalURL retrieves the short URL associated with the given original URL.
func (s *St) GetByOriginalURL(originalURL string) (string, bool) {
	for key, val := range s.URLMap {
		if val == originalURL {
			return key, true
		}
	}
	return "", false
}

// ListAll always returns nil.
func (s *St) ListAll(_ string) ([]*entities.ReqListAll, error) {
	return nil, nil
}

// CreateShortURLs always returns nil.
func (s *St) CreateShortURLs(_ []*entities.ReqURL, _ string) error {
	return nil
}

// DeleteURLs always returns nil.
func (s *St) DeleteURLs(_ []string, _ string) error {
	return nil
}

// URLDeleted always returns false.
func (s *St) URLDeleted(_ string) bool {
	return false
}

// WriteEvent writes an event to the JSON file.
func (s *St) WriteEvent(event *Event) error {
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	defer file.Close()

	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err = writer.Write(data); err != nil {
		return err
	}

	if err = writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()
}

// SyncData synchronizes the in-memory storage data by writing events to the JSON file.
func (s *St) SyncData() {
	for k, v := range s.URLMap {
		event := Event{
			strconv.Itoa(s.lastID),
			k,
			v,
		}

		err := s.WriteEvent(&event)
		if err != nil {
			log.Fatalf("Sync data - %s", err.Error())
		}
	}
}

// Ping pings the in-memory storage (always returns nil).
func (s *St) Ping() error {
	return nil
}
