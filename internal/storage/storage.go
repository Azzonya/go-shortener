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
	URLMap  map[string]string
	file    *os.File
	writer  *bufio.Writer
	maxUUID int
}

type Event struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStorage(filePath string) (*Storage, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	newDecoder := json.NewDecoder(file)
	URLMap := make(map[string]string)

	maxUUID := 0
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
		maxUUID += 1
		URLMap[event.OriginalURL] = event.ShortURL
	}
	return &Storage{
		URLMap:  URLMap,
		file:    file,
		writer:  bufio.NewWriter(file),
		maxUUID: maxUUID,
	}, nil
}

func (s *Storage) Add(key, value string) error {
	s.URLMap[key] = value
	if s.file == nil {
		return nil
	}
	s.maxUUID += 1
	event := Event{
		strconv.Itoa(s.maxUUID),
		key,
		value,
	}
	err := s.WriteEvent(&event)

	return err
}

func (s *Storage) GetOne(key string) (string, bool) {
	URL, exist := s.URLMap[key]
	return URL, exist
}

func (s *Storage) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := s.writer.Write(data); err != nil {
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		return err
	}

	return s.writer.Flush()
}
