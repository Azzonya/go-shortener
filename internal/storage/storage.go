package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

type Storage struct {
	URLMap   map[string]string
	filePath string
	lastID   int
	useDB    bool
}

type Event struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Storage) RestoreFromFile() error {
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

func (s *Storage) SyncData() {
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

func (s *Storage) Add(key, value string) error {
	s.URLMap[key] = value
	s.lastID += 1

	return nil
}

func (s *Storage) GetOne(key string) (string, bool) {
	URL, exist := s.URLMap[key]
	return URL, exist
}

func (s *Storage) WriteEvent(event *Event) error {
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

func NewStorage(filePath string, useDB bool) (*Storage, error) {
	s := &Storage{}

	s.filePath = filePath
	s.useDB = useDB

	if !s.useDB {
		err := s.RestoreFromFile()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
