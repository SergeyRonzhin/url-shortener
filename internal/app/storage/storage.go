package storage

type Storage struct {
	links map[string]string
}

type Repository interface {
	Get(key string) (string, bool)
	Add(key string, value string)
	ContainsValue(value string) (bool, string)
}

func NewStorage() Storage {
	return Storage{make(map[string]string)}
}

func (s Storage) Get(key string) (string, bool) {
	value, exist := s.links[key]
	return value, exist
}

func (s Storage) Add(key string, value string) {
	s.links[key] = value
}

func (s Storage) ContainsValue(value string) (bool, string) {
	for key, url := range s.links {
		if url == value {
			return true, key
		}
	}

	return false, ""
}
