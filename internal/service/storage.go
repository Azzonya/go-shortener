package service

type Storage struct {
	urlMap map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		make(map[string]string),
	}
}

func (s *Storage) add(key, value string) {
	s.urlMap[key] = value
}

func (s *Storage) getOne(key string) (string, bool) {
	URL, exist := s.urlMap[key]
	return URL, exist
}
