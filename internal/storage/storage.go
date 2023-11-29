package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Storage struct {
	URLMap map[string]string
	file   *os.File
	writer *bufio.Writer
	dump   []Event
	lastID int
}

type Event struct {
	NumberUUID  string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Storage) RestoreFromFile() error {
	newDecoder := json.NewDecoder(s.file)

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
	for i := 0; i < len(s.dump); i++ {
		err := s.WriteEvent(&s.dump[i])
		if err != nil {
			log.Fatalf("Sync data - %s", err.Error())
		}

		s.dump = append(s.dump[:i], s.dump[i+1:]...)
		i--
	}
}

func (s *Storage) Add(key, value string) error {
	s.URLMap[key] = value
	if s.file == nil {
		return nil
	}
	s.lastID += 1
	event := Event{
		strconv.Itoa(s.lastID),
		key,
		value,
	}
	s.dump = append(s.dump, event)

	return nil
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
		fmt.Println(err.Error(), 1)
		return err
	}

	if err := s.writer.WriteByte('\n'); err != nil {
		fmt.Println(err.Error(), 2)
		return err
	}

	return s.writer.Flush()
}

func NewStorage(filePath string) (*Storage, error) {
	var err error

	s := &Storage{}

	s.file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	s.writer = bufio.NewWriter(s.file)

	err = s.RestoreFromFile()
	if err != nil {
		return nil, err
	}

	return s, nil
}
