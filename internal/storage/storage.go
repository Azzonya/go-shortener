package storage

type Storage struct {
	UrlMap map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		make(map[string]string),
	}
}

func (s *Storage) Add(key, value string) {
	s.UrlMap[key] = value
}

func (s *Storage) GetOne(key string) (string, bool) {
	URL, exist := s.UrlMap[key]
	return URL, exist
}
