package inmemory

import (
	"bufio"
	"encoding/json"
	"github.com/Azzonya/go-shortener/internal/entities"
	"io"
	"log"
	"os"
	"strconv"
)

type St struct {
	URLMap   map[string]string
	filePath string
	lastID   int
}

type Event struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func New(filePath string) (*St, error) {
	s := &St{}

	s.filePath = filePath
	err := s.Initialize()
	if err != nil {
		return nil, err
	}

	return s, nil
}

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
			if err == io.EOF {
				break
			} else {
				log.Println("error decode JSON:", err)
				break
			}
		}
		s.lastID += 1
		s.URLMap[event.OriginalURL] = event.ShortURL
	}

	return nil
}

func (s *St) TableExist() bool {
	return true
}

func (s *St) Add(originalURL, shortURL, userID string) error {
	s.URLMap[shortURL] = originalURL
	s.lastID += 1

	return nil
}

func (s *St) Update(originalURL, shortURL string) error {
	return nil
}

func (s *St) GetByShortURL(shortURL string) (string, bool) {
	URL, exist := s.URLMap[shortURL]
	return URL, exist
}

func (s *St) GetByOriginalURL(originalURL string) (string, bool) {
	for key, val := range s.URLMap {
		if val == originalURL {
			return key, true
		}
	}
	return "", false
}

func (s *St) ListAll(userID string) ([]*entities.ReqListAll, error) {
	return nil, nil
}

func (s *St) CreateShortURLs(urls []*entities.ReqURL, userID string) error {
	return nil
}

func (s *St) DeleteURLs(urls []string, userID string) error {
	return nil
}

func (s *St) URLDeleted(shortURL string) bool {
	return false
}

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

func (s *St) Ping() error {
	return nil
}
