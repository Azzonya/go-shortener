package storage

type Storage struct {
	URLMap map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		make(map[string]string),
	}
}

func (s *Storage) Add(key, value string) {
	s.URLMap[key] = value
}

func (s *Storage) GetOne(key string) (string, bool) {
	URL, exist := s.URLMap[key]
	return URL, exist
}
