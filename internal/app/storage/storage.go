package storage

type Storage struct {
	links map[string]string
}

func NewStorage() Storage {
	return Storage{make(map[string]string)}
}

func (s Storage) Get(key string) string {
	return s.links[key]
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
